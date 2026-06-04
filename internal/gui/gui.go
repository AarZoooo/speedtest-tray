package gui

import (
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"speedtest-tray/internal/speedtest_util"
)

// transparentTheme makes the window background invisible to allow rounded corners
type transparentTheme struct {
	fyne.Theme
}

func (t *transparentTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameBackground {
		return color.Transparent
	}
	return t.Theme.Color(name, variant)
}

type GUI struct {
	App    fyne.App
	Window fyne.Window

	// UI Elements
	titleLabel  *widget.Label
	statusLabel *widget.Label // Visible status above start button
	serverLabel *widget.Label
	pingLabel   *widget.Label
	dlLabel     *widget.Label
	dlProgress  *narrowProgressBar
	ulLabel     *widget.Label
	ulProgress  *narrowProgressBar
	runButton   *widget.Button

	tester *speedtest_util.SpeedTester

	// Animation state
	spinnerStop chan struct{}
}

func New(tester *speedtest_util.SpeedTester) *GUI {
	log.Println("Initializing GUI...")
	myApp := app.NewWithID("com.minispeedtest.app")

	// Apply transparency theme
	myApp.Settings().SetTheme(&transparentTheme{Theme: theme.DefaultTheme()})

	var myWindow fyne.Window
	if desk, ok := myApp.Driver().(desktop.Driver); ok {
		myWindow = desk.CreateSplashWindow()
	} else {
		myWindow = myApp.NewWindow("Speedtest Tray")
	}

	myWindow.SetTitle("Speedtest Tray")

	g := &GUI{
		App:    myApp,
		Window: myWindow,
		tester: tester,
	}

	g.initUI()
	g.setupTray()

	myWindow.Resize(fyne.NewSize(320, 460))
	myWindow.SetFixedSize(true)

	return g
}

func (g *GUI) initUI() {
	g.titleLabel = widget.NewLabel("Speedtest Tray")
	g.titleLabel.Alignment = fyne.TextAlignCenter
	g.titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	g.titleLabel.SizeName = theme.SizeNameHeadingText

	closeBtn := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		g.Window.Hide()
	})
	closeBtn.Importance = widget.LowImportance

	header := container.NewBorder(nil, nil, nil, closeBtn, g.titleLabel)

	g.statusLabel = widget.NewLabel("Ready")
	g.statusLabel.Alignment = fyne.TextAlignCenter
	g.statusLabel.Importance = widget.LowImportance

	// Metrics
	serverTitle := widget.NewLabel("Server Location:")
	g.serverLabel = widget.NewLabel("--")
	g.serverLabel.Alignment = fyne.TextAlignTrailing
	g.serverLabel.TextStyle = fyne.TextStyle{Bold: true}

	pingTitle := widget.NewLabel("Ping:")
	g.pingLabel = widget.NewLabel("--")
	g.pingLabel.Alignment = fyne.TextAlignTrailing
	g.pingLabel.TextStyle = fyne.TextStyle{Bold: true}

	dlTitle := widget.NewLabel("Download:")
	g.dlLabel = widget.NewLabel("--")
	g.dlLabel.Alignment = fyne.TextAlignTrailing
	g.dlLabel.TextStyle = fyne.TextStyle{Bold: true}

	ulTitle := widget.NewLabel("Upload:")
	g.ulLabel = widget.NewLabel("--")
	g.ulLabel.Alignment = fyne.TextAlignTrailing
	g.ulLabel.TextStyle = fyne.TextStyle{Bold: true}

	g.dlProgress = newNarrowProgressBar()
	g.ulProgress = newNarrowProgressBar()

	g.runButton = widget.NewButton("Start", func() {
		g.startTest()
	})

	metricsGrid := container.NewGridWithColumns(2,
		serverTitle, g.serverLabel,
		pingTitle, g.pingLabel,
		dlTitle, g.dlLabel,
	)

	ulRow := container.NewGridWithColumns(2, ulTitle, g.ulLabel)

	metrics := container.NewVBox(
		metricsGrid,
		g.dlProgress,
		ulRow,
		g.ulProgress,
	)

	// Inner content with padding
	innerContent := container.NewVBox(
		header,
		container.NewPadded(metrics),
	)

	bottomContent := container.NewVBox(
		container.NewPadded(g.statusLabel),
		g.runButton,
	)

	mainUI := container.NewBorder(innerContent, bottomContent, nil, nil)

	// Create rounded background
	bgColor := theme.DefaultTheme().Color(theme.ColorNameBackground, theme.VariantDark)
	// If variant is light, we'd need to check, but let's stick to dark-ish default or detect
	bg := canvas.NewRectangle(bgColor)
	bg.CornerRadius = 16

	// Wrap the entire UI in a Stack with the rounded background and overall padding
	paddedUI := container.NewPadded(container.NewPadded(mainUI))

	finalLayout := container.NewStack(bg, paddedUI)

	g.Window.SetContent(finalLayout)
}
