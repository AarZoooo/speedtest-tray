package config

import "time"

// App branding
const AppName = "SpeedTest Tray"

// Progress thresholds (phase orchestration)
const (
	ProgressInit         = 0.05
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
	TestDurationDownload = 10 * time.Second
	TestDurationUpload   = 10 * time.Second
	ResultTimeout        = 10 * time.Second
)

// Error messages
const (
	ErrTestStopped = "test stopped"
	ErrTestTimeout = "test timeout"
)

// UI timing
const (
	UIHideDelayMs = 1000
)

// Window properties
const (
	WindowWidth         = 320
	WindowHeight        = 560
	WindowCornerRadius  = 32
	WindowOffsetYPixels = -10
)

// Gauge scales (Mbps)
const (
	GaugeMaxDownload = 1000
	GaugeMaxUpload   = 100
)
