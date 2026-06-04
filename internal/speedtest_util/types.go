package speedtest_util

import "speedtest-tray/internal/config"

// Phase alias for easier use
type Phase = config.Phase

// Phase constants
const (
	INITIALIZING      = config.PhaseInitializing
	GETTING_INFO      = config.PhaseGettingInfo
	FINDING_SERVERS   = config.PhaseFindingServers
	SELECTING_SERVER  = config.PhaseSelectingServer
	SERVER_SELECTED   = config.PhaseServerSelected
	PING_TEST         = config.PhasePingTest
	STARTING_DOWNLOAD = config.PhaseStartingDownload
	DOWNLOADING       = config.PhaseDownloading
	STARTING_UPLOAD   = config.PhaseStartingUpload
	UPLOADING         = config.PhaseUploading
	COMPLETED         = config.PhaseCompleted
	FAILED            = config.PhaseFailed
)

type Update struct {
	Phase    Phase
	Progress float64
	Ping     float64
	Download float64
	Upload   float64
	Server   string
	Error    error
}

type Result struct {
	Phase    Phase
	Ping     float64
	Download float64
	Upload   float64
	Server   string
	Error    error
}
