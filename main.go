package main

//go:generate go run ./cmd/gen-frontend-config

import (
	"bytes"
	"context"
	"embed"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"speedtest-tray/internal/cli"
	"speedtest-tray/internal/config"
	"speedtest-tray/internal/gui_wails"
	"speedtest-tray/internal/speedtest_util"
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
	cliFlag := flag.Bool(config.FlagCLI, false, config.UsageCLI)
	cliFlagShort := flag.Bool(config.FlagCLIShort, false, config.UsageCLI)
	jsonFlag := flag.Bool(config.FlagJSON, false, config.UsageJSON)
	jsonFlagShort := flag.Bool(config.FlagJSONShort, false, config.UsageJSON)
	serverFlag := flag.String(config.FlagServer, "", config.UsageServer)
	serverFlagShort := flag.String(config.FlagServerShort, "", config.UsageServer)
	historyFlag := flag.Bool(config.FlagHistory, false, config.UsageHistory)
	historyFlagShort := flag.Bool(config.FlagHistoryShort, false, config.UsageHistory)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", config.AppName)
		fmt.Fprintf(os.Stderr, "  %s [flags]\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Flags:")
		flag.PrintDefaults()
	}
	flag.Parse()

	isCLI := *cliFlag || *cliFlagShort || *jsonFlag || *jsonFlagShort || *serverFlag != "" || *serverFlagShort != "" || *historyFlag || *historyFlagShort
	if isCLI {
		attachConsole()
		jsonMode := *jsonFlag || *jsonFlagShort

		if *historyFlag || *historyFlagShort {
			if err := cli.PrintHistory(os.Stdout, jsonMode); err != nil {
				os.Exit(1)
			}
			os.Exit(0)
		}

		serverID := *serverFlag
		if serverID == "" {
			serverID = *serverFlagShort
		}

		appConfig = config.LoadConfigOrDefault()
		if appConfig.SaveLogs {
			enableFileLogging()
		} else {
			slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
		}

		tester := speedtest_util.New()
		tester.TargetServerID = serverID

		ctx := context.Background()
		err := cli.Run(ctx, jsonMode, serverID, tester, os.Stdout)
		if err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}

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
	}, func(enabled bool) {
		app.SetLaunchAtLogin(enabled)
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
			WindowIsTranslucent:               true,
			BackdropType:                      windows.Acrylic,
			DisableWindowIcon:                 false,
			DisableFramelessWindowDecorations: false,
			IsZoomControlEnabled:              false,
			DisablePinchZoom:                  true,
		},
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHidden(),
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
		},
	}
}
