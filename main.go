package main

//go:generate go run ./cmd/gen-frontend-config

import (
	"bytes"
	"embed"
	"log/slog"
	"os"

	"github.com/energye/systray"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"speedtest-tray/internal/config"
	"speedtest-tray/internal/gui_wails"
)

//go:embed all:frontend
var assets embed.FS

//go:embed build/windows/icon.ico
var iconBytes []byte

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

	startTray(app)

	options := newOptions(app)
	options.Logger = logger.NewDefaultLogger()

	if err := wails.Run(options); err != nil {
		appLogger.Error(config.ErrRunWails, "error", err)
		os.Exit(1)
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

		saveLogs := systray.AddMenuItemCheckbox("Enable Session Logging", "Save test logs to app data", appConfig.SaveLogs)
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
			config.SaveConfig(appConfig)
		})

		openLogs := systray.AddMenuItem("Open Logs Directory", "Open the directory containing the session logs")
		openLogs.Click(func() {
			if err := config.OpenDirectory(config.GetConfigDir()); err != nil {
				appLogger.Error(config.ErrOpenLogsDir, "error", err)
			}
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
	logDir := config.GetConfigDir()

	if err := os.MkdirAll(logDir, 0755); err != nil {
		appLogger.Error(config.ErrCreateLogDir, "error", err)
		return
	}

	logPath := config.GetLogFilePath()
	truncateLogFile(logPath, config.MaxLogLines)

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

func truncateLogFile(logPath string, maxLines int) {
	data, err := os.ReadFile(logPath)
	if err != nil {
		return
	}

	lines := bytes.Split(data, []byte("\n"))
	if len(lines) > 0 && len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}

	if len(lines) <= maxLines {
		return
	}

	trimmedLines := lines[len(lines)-maxLines:]
	output := bytes.Join(trimmedLines, []byte("\n"))
	output = append(output, '\n')

	_ = os.WriteFile(logPath, output, 0666)
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
