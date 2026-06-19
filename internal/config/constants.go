package config

import "time"

// App branding
const AppName = "SpeedTest Tray"

const AppVersion = "1.1.1"

const (
	GitHubOwner = "AarZoooo"
	GitHubRepo  = "speedtest-tray"
)

// Progress thresholds (phase orchestration)
const (
	ProgressInit         = 0.00
	ProgressGetInfo      = 0.05
	ProgressFindServers  = 0.10
	ProgressSelectServer = 0.12
	ProgressServerSelect = 0.15
	ProgressPingStart    = 0.20
	ProgressPingEnd      = 0.30
	ProgressDownStart    = 0.30
	ProgressDownEnd      = 0.70
	ProgressUpStart      = 0.70
	ProgressUpEnd        = 0.95
	ProgressComplete     = 1.0
)

// Test durations
const (
	EstimatedDurationDownload = 15 * time.Second
	EstimatedDurationUpload   = 15 * time.Second
	ResultTimeout             = 10 * time.Second
)

// Retry configuration
const (
	MaxRetryAttempts = 3
	RetryDelay       = 1 * time.Second
)

// Error messages
const (
	ErrTestStopped = "test stopped"
	ErrTestTimeout = "test timeout"
	ErrNoInternet  = "no internet connection"
	MsgNoInternet  = "No internet connection"
)

// Connectivity probe
const (
	ConnectivityProbeURL     = "https://clients3.google.com/generate_204"
	ConnectivityCheckTimeout = 5 * time.Second
)

// UI timing
const (
	UIHideDelayMs      = 2000
	PhaseSleepDuration = 2 * time.Second
	ToggleThreshold    = 200 * time.Millisecond
)

// Window properties
const (
	WindowWidth         = 320
	WindowHeight        = 560
	WindowCornerRadius  = 16
	WindowOffsetYPixels = -20
	StandardDPI         = 96.0
)

// Win32 constants
const (
	MonitorDefaultToNearest = 2
	MdtEffectiveDpi         = 0
	SystrayIconID           = 100
	WindowTrayGap           = 8
)

// Gauge scales (Mbps)
const (
	GaugeMaxDownload = 1000
	GaugeMaxUpload   = 100
)

// Log messages
const (
	LogAppStarting     = "--- Application Starting ---"
	LogLoggingEnabled  = "--- File Logging Enabled ---"
	LogLoggingDisabled = "--- File Logging Disabled ---"
	ErrRunWails        = "Failed to run Wails app"
	ErrCreateLogDir    = "Failed to create log directory"
	ErrOpenLogFile     = "Failed to open log file"
	ErrCreateConfigDir = "Failed to create config directory"
	ErrOpenLogsDir     = "Failed to open logs directory"
	LogAdapterUpdate   = "Adapter: Update"
	LogAdapterClosed   = "Adapter: Updates closed, waiting for result"
	LogAdapterResult   = "Adapter: Result received"
	LogAdapterTimeout  = "Adapter: Timeout waiting for result"
	LogHardwareStats   = "Hardware Utilization"

	LogUpdateCheckStart  = "Update check started"
	LogUpdateFound       = "Update available"
	LogUpdateNoneFound   = "No update available"
	LogUpdateApplying    = "Applying update"
	LogUpdateCleanup     = "Cleaned up staged installer"
	ErrUpdateCheck       = "Failed to check for update"
	ErrUpdateApply       = "Failed to apply update"
	ErrUpdateDownload    = "Failed to download update"
	ErrUpdateBadChecksum = "Downloaded installer size mismatch, aborting"
	ErrUpdateSkip        = "Failed to skip update version"
	ErrAutostartEnable   = "Failed to enable launch at login"
	ErrAutostartDisable  = "Failed to disable launch at login"
)

const (
	KeyAllocMB      = "alloc_mb"
	KeySysMB        = "sys_mb"
	KeyNumGoroutine = "num_goroutine"
	KeyPhase        = "phase"
	KeyProgress     = "progress"
	KeyError        = "error"
)

const (
	UpdateChannelSize = 64
	MaxLogLines       = 5000
	MaxHistoryEntries = 50
)

const (
	FlagCLI            = "cli"
	FlagCLIShort       = "c"
	FlagJSON           = "json"
	FlagJSONShort      = "j"
	FlagServer         = "server"
	FlagServerShort    = "s"
	FlagHistory        = "history"
	FlagHistoryShort   = "h"
	UsageCLI           = "Run in headless CLI mode"
	UsageJSON          = "Output results as JSON (implies CLI mode)"
	UsageServer        = "Target server ID (implies CLI mode)"
	UsageHistory       = "Show speedtest history (implies CLI mode)"
	CLIHeader          = "Speedtest Tray (Headless CLI)"
	CLILineSeparator   = "------------------------------------------------------------"
	CLIDoubleLine      = "============================================================"
	CLIResultHeader    = "                     TEST RESULTS"
	CLIStatusChecking  = "Checking connectivity..."
	CLIStatusSelecting = "Selecting server..."
	CLIStatusPing      = "Running Ping test..."
	CLIStatusDownload  = "Running Download test..."
	CLIStatusUpload    = "Running Upload test..."
	CLIStatusCompleted = "Done."
	JSONStatusSuccess  = "success"
	JSONStatusFailed   = "failed"
)
