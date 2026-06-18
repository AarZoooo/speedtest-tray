package speedtest_util

import (
	"github.com/showwin/speedtest-go/speedtest"
)

// SpeedTester implements TestOrchestrator
type SpeedTester struct {
	client         *speedtest.Speedtest
	user           *speedtest.User
	servers        speedtest.Servers
	server         *speedtest.Server
	TargetServerID string
}


// New creates a new SpeedTester instance
func New() *SpeedTester {
	return &SpeedTester{
		client: speedtest.New(),
	}
}
