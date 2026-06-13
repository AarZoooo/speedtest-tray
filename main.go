package main

//go:generate go run ./cmd/gen-frontend-config

import (
	"embed"
	"log/slog"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"speedtest-tray/internal/config"
	"speedtest-tray/internal/gui_wails"
)

//go:embed all:frontend
var assets embed.FS

//go:embed build/windows/icon.ico
var iconBytes []byte

//go:embed build/darwin/iconTemplate.png
var macIconBytes []byte

var (
	logFile          *os.File
	isLoggingEnabled bool
	appConfig        config.CustomConfig
	appLogger        *slog.Logger
)

func main() {
	appLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(appLogger)

	appLogger.Info(config.LogAppStarting)

	app := gui_wails.NewApp()

	appConfig = config.LoadConfigOrDefault()
	if appConfig.SaveLogs {
		enableFileLogging()
	}

	gui_wails.StartTray(app, iconBytes, macIconBytes, &appConfig, func(enabled bool) {
		if enabled {
			appConfig.SaveLogs = true
			enableFileLogging()
		} else {
			appConfig.SaveLogs = false
			disableFileLogging()
		}
		config.SaveConfig(appConfig)
	})

	options := newOptions(app)
	options.Logger = logger.NewDefaultLogger()

	if err := wails.Run(options); err != nil {
		appLogger.Error(config.ErrRunWails, "error", err)
		os.Exit(1)
	}
}

func enableFileLogging() {
	logDir := config.GetConfigDir()

	if err := os.MkdirAll(logDir, 0755); err != nil {
		appLogger.Error(config.ErrCreateLogDir, "error", err)
		return
	}

	logPath := config.GetLogFilePath()
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		appLogger.Error(config.ErrOpenLogFile, "error", err)
		return
	}

	logFile = file
	handler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: slog.LevelDebug})
	appLogger = slog.New(handler)
	slog.SetDefault(appLogger)
	isLoggingEnabled = true
	appLogger.Info(config.LogLoggingEnabled)
}

func disableFileLogging() {
	if logFile != nil {
		appLogger.Info(config.LogLoggingDisabled)
		appLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		slog.SetDefault(appLogger)
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
		OnShutdown:       app.Shutdown,
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
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHidden(),
			WebviewIsTransparent: true,
			WindowIsTranslucent:  false,
		},
	}
}
