package main

import (
	"embed"
	"log"
	"os"

	"github.com/energye/systray"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"speedtest-tray/internal/config"
	"speedtest-tray/internal/gui_wails"
	"speedtest-tray/internal/speedtest_util"
)

//go:embed all:frontend
var assets embed.FS

//go:embed build/windows/icon.ico
var iconBytes []byte

var (
	logFile          *os.File
	isLoggingEnabled bool
	appConfig        config.AppConfig
)

func main() {
	log.Println("--- Application Starting ---")

	tester := speedtest_util.New()
	app := gui_wails.NewApp(tester)

	appConfig = config.Load()
	if appConfig.SaveLogs {
		enableFileLogging()
	}

	startTray(app)

	options := newOptions(app)
	options.Logger = logger.NewDefaultLogger()

	if err := wails.Run(options); err != nil {
		log.Fatal(err)
	}
}

func startTray(app *gui_wails.App) {
	go systray.Run(func() {
		systray.SetTitle(config.AppName)
		systray.SetTooltip(config.AppName)
		systray.SetIcon(iconBytes)
		systray.SetOnClick(func(menu systray.IMenu) {
			go app.ShowWindow()
		})

		show := systray.AddMenuItem("Show", "Show the speedtest window")
		show.Click(app.ShowWindow)

		systray.AddSeparator()

		saveLogs := systray.AddMenuItemCheckbox("Save logs to Documents", "Enable or disable file logging", appConfig.SaveLogs)
		saveLogs.Click(func() {
			if saveLogs.Checked() {
				saveLogs.Uncheck()
				appConfig.SaveLogs = false
				disableFileLogging()
			} else {
				saveLogs.Check()
				appConfig.SaveLogs = true
				enableFileLogging()
			}
			config.Save(appConfig)
		})

		systray.AddSeparator()

		quit := systray.AddMenuItem("Quit", "Quit the application")
		quit.Click(func() {
			app.Quit()
			systray.Quit()
		})
	}, func() {})
}

func enableFileLogging() {
	logDir := config.GetAppDir()

	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Printf("Failed to create log directory: %v\n", err)
		return
	}

	logPath := config.GetLogFilePath()
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Printf("Failed to open log file: %v\n", err)
		return
	}

	logFile = file
	log.SetOutput(logFile)
	isLoggingEnabled = true
	log.Println("--- File Logging Enabled ---")
}

func disableFileLogging() {
	if logFile != nil {
		log.Println("--- File Logging Disabled ---")
		log.SetOutput(os.Stdout)
		logFile.Close()
		logFile = nil
	}
	isLoggingEnabled = false
}

func newOptions(app *gui_wails.App) *options.App {
	return &options.App{
		Title:            config.AppName,
		Width:            config.WindowWidth,
		Height:           config.WindowHeight,
		MinWidth:         config.WindowWidth,
		MinHeight:        config.WindowHeight,
		MaxWidth:         config.WindowWidth,
		MaxHeight:        config.WindowHeight,
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
			DisableWindowIcon:                 false,
			DisableFramelessWindowDecorations: true,
			IsZoomControlEnabled:              false,
			DisablePinchZoom:                  true,
		},
	}
}
