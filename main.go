package main

import (
	"embed"
	"log"

	"github.com/energye/systray"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"speedtest-tray/internal/gui_wails"
	"speedtest-tray/internal/speedtest_util"
)

//go:embed all:frontend
var assets embed.FS

//go:embed build/windows/icon.ico
var iconBytes []byte

func main() {
	tester := speedtest_util.New()
	app := gui_wails.NewApp(tester)

	startTray(app)

	if err := wails.Run(newOptions(app)); err != nil {
		log.Fatal(err)
	}
}

func startTray(app *gui_wails.App) {
	go systray.Run(func() {
		systray.SetTitle(speedtest_util.AppName)
		systray.SetTooltip(speedtest_util.AppName)
		systray.SetIcon(iconBytes)
		systray.SetOnClick(func(menu systray.IMenu) {
			app.ShowWindow()
		})

		show := systray.AddMenuItem("Show", "Show the speedtest window")
		show.Click(app.ShowWindow)

		systray.AddSeparator()

		quit := systray.AddMenuItem("Quit", "Quit the application")
		quit.Click(func() {
			app.Quit()
			systray.Quit()
		})
	}, func() {})
}

func newOptions(app *gui_wails.App) *options.App {
	return &options.App{
		Title:            speedtest_util.AppName,
		Width:            gui_wails.WindowWidth,
		Height:           gui_wails.WindowHeight,
		MinWidth:         gui_wails.WindowWidth,
		MinHeight:        gui_wails.WindowHeight,
		MaxWidth:         gui_wails.WindowWidth,
		MaxHeight:        gui_wails.WindowHeight,
		DisableResize:    true,
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		OnStartup:        app.Startup,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Bind: []interface{}{
			app,
		},
		Frameless:         true,
		AlwaysOnTop:       true,
		StartHidden:       true,
		HideWindowOnClose: true,
		Windows: &windows.Options{
			WebviewIsTransparent:              true,
			WindowIsTranslucent:               false,
			BackdropType:                      windows.None,
			DisableWindowIcon:                 true,
			DisableFramelessWindowDecorations: true,
			IsZoomControlEnabled:              false,
			DisablePinchZoom:                  true,
		},
	}
}
