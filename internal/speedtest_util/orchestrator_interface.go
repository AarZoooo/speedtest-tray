package speedtest_util

import (
	"context"
	"time"
)

// TestOrchestrator defines the contract for running speed tests.
// This allows for mocking and testing without external dependencies.
type TestOrchestrator interface {
	// GetUserInfo fetches user location and ISP info
	GetUserInfo(ctx context.Context) error

	// FindServers fetches available speed test servers
	FindServers(ctx context.Context) error

	// SelectBestServer selects the closest server
	SelectBestServer(ctx context.Context) (*ServerInfo, error)

	// RunPing executes a ping test
	RunPing(ctx context.Context) (time.Duration, error)

	// RunDownload executes a download test with progress callback
	RunDownload(ctx context.Context, callback func(float64)) (float64, error)

	// RunUpload executes an upload test with progress callback
	RunUpload(ctx context.Context, callback func(float64)) (float64, error)
}

// ServerInfo holds selected server details
type ServerInfo struct {
	Name    string
	Country string
}
