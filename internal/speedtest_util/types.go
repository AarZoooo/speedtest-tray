package speedtest_util

const AppName = "SpeedTest Tray"

type Phase string

const (
	INITIALIZING      Phase = "INITIALIZING"
	GETTING_INFO      Phase = "GETTING_INFO"
	FINDING_SERVERS   Phase = "FINDING_SERVERS"
	SELECTING_SERVER  Phase = "SELECTING_SERVER"
	SERVER_SELECTED   Phase = "SERVER_SELECTED"
	PING_TEST         Phase = "PING_TEST"
	STARTING_DOWNLOAD Phase = "STARTING_DOWNLOAD"
	DOWNLOADING       Phase = "DOWNLOADING"
	STARTING_UPLOAD   Phase = "STARTING_UPLOAD"
	UPLOADING         Phase = "UPLOADING"
	COMPLETED         Phase = "COMPLETED"
	FAILED            Phase = "FAILED"
)

func (p Phase) String() string {
	switch p {
	case INITIALIZING:
		return "Initializing..."
	case GETTING_INFO:
		return "Fetching user info..."
	case FINDING_SERVERS:
		return "Locating servers..."
	case SELECTING_SERVER:
		return "Selecting best server..."
	case SERVER_SELECTED:
		return "Server selected"
	case PING_TEST:
		return "Measuring latency..."
	case STARTING_DOWNLOAD:
		return "Starting download test..."
	case DOWNLOADING:
		return "Running download test..."
	case STARTING_UPLOAD:
		return "Starting upload test..."
	case UPLOADING:
		return "Running upload test..."
	case COMPLETED:
		return "Test Completed"
	case FAILED:
		return "Test Failed"
	default:
		return string(p)
	}
}

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
