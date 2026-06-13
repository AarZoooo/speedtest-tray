package config

import "time"

// App branding
const AppName = "SpeedTest Tray"

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
	ErrTestStopped  = "test stopped"
	ErrTestTimeout  = "test timeout"
	ErrNoInternet   = "no internet connection"
	MsgNoInternet   = "No internet connection"
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
)

// Window properties
const (
	WindowWidth         = 320
	WindowHeight        = 560
	WindowCornerRadius  = 32
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
	LogAppStarting    = "--- Application Starting ---"
	LogLoggingEnabled  = "--- File Logging Enabled ---"
	LogLoggingDisabled = "--- File Logging Disabled ---"
	ErrRunWails       = "Failed to run Wails app"
	ErrCreateLogDir   = "Failed to create log directory"
	ErrOpenLogFile    = "Failed to open log file"
	ErrCreateConfigDir = "Failed to create config directory"
	LogAdapterUpdate     = "Adapter: Update"
	LogAdapterClosed     = "Adapter: Updates closed, waiting for result"
	LogAdapterResult     = "Adapter: Result received"
	LogAdapterTimeout    = "Adapter: Timeout waiting for result"
	LogHardwareStats     = "Hardware Utilization"
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
)
