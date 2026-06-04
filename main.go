package main

import (
	"embed"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/energye/systray"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
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

type Config struct {
	SaveLogs bool `json:"save_logs"`
}

var (
	logFile          *os.File
	isLoggingEnabled bool
	config           Config
)

func main() {
	// Start with standard logging
	log.Println("--- Application Starting ---")

	tester := speedtest_util.New()
	app := gui_wails.NewApp(tester)

	// Load config
	loadConfig()

	startTray(app)

	options := newOptions(app)
	options.Logger = logger.NewDefaultLogger()

	if err := wails.Run(options); err != nil {
		log.Fatal(err)
	}
}

func getAppDir() string {
	home, _ := os.UserHomeDir()

	// Try OneDrive first
	dir := filepath.Join(home, "OneDrive", "Documents", "SpeedTest Tray")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		dir = filepath.Join(home, "Documents", "SpeedTest Tray")
	}
	return dir
}

func loadConfig() {
	dir := getAppDir()
	configPath := filepath.Join(dir, "config.json")

	data, err := os.ReadFile(configPath)
	if err == nil {
		json.Unmarshal(data, &config)
	}

	if config.SaveLogs {
		enableFileLogging()
	}
}

func saveConfig() {
	dir := getAppDir()
	os.MkdirAll(dir, 0755)
	configPath := filepath.Join(dir, "config.json")

	data, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile(configPath, data, 0644)
}

func startTray(app *gui_wails.App) {
	go systray.Run(func() {
		systray.SetTitle(speedtest_util.AppName)
		systray.SetTooltip(speedtest_util.AppName)
		systray.SetIcon(iconBytes)
		systray.SetOnClick(func(menu systray.IMenu) {
			go app.ShowWindow()
		})

		show := systray.AddMenuItem("Show", "Show the speedtest window")
		show.Click(app.ShowWindow)

		systray.AddSeparator()

		saveLogs := systray.AddMenuItemCheckbox("Save logs to Documents", "Enable or disable file logging", config.SaveLogs)
		saveLogs.Click(func() {
			if saveLogs.Checked() {
				saveLogs.Uncheck()
				config.SaveLogs = false
				disableFileLogging()
			} else {
				saveLogs.Check()
				config.SaveLogs = true
				enableFileLogging()
			}
			saveConfig()
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
	logDir := getAppDir()

	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Printf("Failed to create log directory: %v\n", err)
		return
	}

	logPath := filepath.Join(logDir, "app.log")
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Printf("Failed to open log file: %v\n", err)
		return
	}

	logFile = file
	// IMPORTANT: On Windows, logging to os.Stdout in a GUI app can block if not consumed.
	// We'll log ONLY to the file when enabled to be safe.
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
			DisableWindowIcon:                 false,
			DisableFramelessWindowDecorations: true,
			IsZoomControlEnabled:              false,
			DisablePinchZoom:                  true,
		},
	}
}
