package config

// Phase represents a stage in the speed test lifecycle
type Phase string

const (
	PhaseInitializing      Phase = "INITIALIZING"
	PhaseGettingInfo       Phase = "GETTING_INFO"
	PhaseFindingServers    Phase = "FINDING_SERVERS"
	PhaseSelectingServer   Phase = "SELECTING_SERVER"
	PhaseServerSelected    Phase = "SERVER_SELECTED"
	PhasePingTest          Phase = "PING_TEST"
	PhaseStartingDownload  Phase = "STARTING_DOWNLOAD"
	PhaseDownloading       Phase = "DOWNLOADING"
	PhaseStartingUpload    Phase = "STARTING_UPLOAD"
	PhaseUploading         Phase = "UPLOADING"
	PhaseCompleted         Phase = "COMPLETED"
	PhaseFailed            Phase = "FAILED"
)

// String returns human-readable phase name
func (p Phase) String() string {
	switch p {
	case PhaseInitializing:
		return "Initializing..."
	case PhaseGettingInfo:
		return "Fetching user info..."
	case PhaseFindingServers:
		return "Locating servers..."
	case PhaseSelectingServer:
		return "Selecting best server..."
	case PhaseServerSelected:
		return "Server selected"
	case PhasePingTest:
		return "Measuring latency..."
	case PhaseStartingDownload:
		return "Starting download test..."
	case PhaseDownloading:
		return "Running download test..."
	case PhaseStartingUpload:
		return "Starting upload test..."
	case PhaseUploading:
		return "Running upload test..."
	case PhaseCompleted:
		return "Test Completed"
	case PhaseFailed:
		return "Test Failed"
	default:
		return string(p)
	}
}
