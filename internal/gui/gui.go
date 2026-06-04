package gui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"speedtest-tray/internal/speedtest_util"
)

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
	myApp := app.New()
	myWindow := myApp.NewWindow("Speedtest-Tray")

	g := &GUI{
		App:    myApp,
		Window: myWindow,
		tester: tester,
	}

	g.initUI()

	// Set window to a small, fixed size typical for a taskbar pop-up
	myWindow.Resize(fyne.NewSize(300, 420))
	myWindow.SetFixedSize(true)

	return g
}

func (g *GUI) initUI() {
	g.titleLabel = widget.NewLabel("SpeedTest Tray")
	g.titleLabel.Alignment = fyne.TextAlignCenter
	g.titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	g.titleLabel.SizeName = theme.SizeNameHeadingText

	g.statusLabel = widget.NewLabel("Ready")
	g.statusLabel.Alignment = fyne.TextAlignCenter
	g.statusLabel.Importance = widget.LowImportance

	// Metrics Labels (Left) and Values (Right)
	serverTitle := widget.NewLabel("Server Location:")
	g.serverLabel = widget.NewLabel("--")
	g.serverLabel.Alignment = fyne.TextAlignTrailing

	pingTitle := widget.NewLabel("Ping:")
	g.pingLabel = widget.NewLabel("--")
	g.pingLabel.Alignment = fyne.TextAlignTrailing

	dlTitle := widget.NewLabel("Download:")
	g.dlLabel = widget.NewLabel("--")
	g.dlLabel.Alignment = fyne.TextAlignTrailing

	ulTitle := widget.NewLabel("Upload:")
	g.ulLabel = widget.NewLabel("--")
	g.ulLabel.Alignment = fyne.TextAlignTrailing

	g.dlProgress = newNarrowProgressBar()
	g.ulProgress = newNarrowProgressBar()

	g.runButton = widget.NewButton("Start", func() {
		g.startTest()
	})

	// Metrics layout: Grid with 2 columns
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

	// Main content in the center
	topContent := container.NewVBox(
		g.titleLabel,
		container.NewPadded(metrics),
	)

	// Use BorderLayout to pin status and button to the bottom
	bottomContent := container.NewVBox(
		container.NewPadded(g.statusLabel),
		g.runButton,
	)
	content := container.NewBorder(topContent, bottomContent, nil, nil)

	g.Window.SetContent(content)
}
