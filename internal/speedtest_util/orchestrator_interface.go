package speedtest_util

import (
	"context"
	"time"
)

type TestOrchestrator interface {
	GetUserInfo(ctx context.Context) error
	FindServers(ctx context.Context) error
	SelectBestServer(ctx context.Context) (*ServerInfo, error)
	RunPing(ctx context.Context) (time.Duration, error)
	RunDownload(ctx context.Context, callback func(float64)) (float64, error)
	RunUpload(ctx context.Context, callback func(float64)) (float64, error)
}

type ServerInfo struct {
	Name    string
	Country string
}
