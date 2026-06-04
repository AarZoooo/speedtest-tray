package speedtest_util

import (
	"context"
	"fmt"
	"time"

	"github.com/showwin/speedtest-go/speedtest"
)

// GetUserInfo implements TestOrchestrator.GetUserInfo
func (st *SpeedTester) GetUserInfo(ctx context.Context) error {
	user, err := st.client.FetchUserInfoContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch user info: %w", err)
	}
	st.user = user
	return nil
}

// FindServers implements TestOrchestrator.FindServers
func (st *SpeedTester) FindServers(ctx context.Context) error {
	servers, err := st.client.FetchServers()
	if err != nil {
		return fmt.Errorf("failed to fetch servers: %w", err)
	}
	st.servers = servers
	return nil
}

// SelectBestServer implements TestOrchestrator.SelectBestServer
func (st *SpeedTester) SelectBestServer(ctx context.Context) (*ServerInfo, error) {
	targets, err := st.servers.FindServer([]int{})
	if err != nil {
		return nil, fmt.Errorf("failed to find target server: %w", err)
	}
	st.server = targets[0]
	return &ServerInfo{
		Name:    st.server.Name,
		Country: st.server.Country,
	}, nil
}

// RunPing implements TestOrchestrator.RunPing
func (st *SpeedTester) RunPing(ctx context.Context) (time.Duration, error) {
	if err := st.server.PingTestContext(ctx, nil); err != nil {
		return 0, fmt.Errorf("ping test failed: %w", err)
	}
	return st.server.Latency, nil
}

// RunDownload implements TestOrchestrator.RunDownload
func (st *SpeedTester) RunDownload(ctx context.Context, callback func(float64)) (float64, error) {
	st.client.SetCallbackDownload(func(rate speedtest.ByteRate) {
		callback(rate.Mbps())
	})

	if err := st.server.DownloadTestContext(ctx); err != nil {
		st.client.SetCallbackDownload(nil)
		return 0, fmt.Errorf("download test failed: %w", err)
	}

	st.client.SetCallbackDownload(nil)
	return st.server.DLSpeed.Mbps(), nil
}

// RunUpload implements TestOrchestrator.RunUpload
func (st *SpeedTester) RunUpload(ctx context.Context, callback func(float64)) (float64, error) {
	st.client.SetCallbackUpload(func(rate speedtest.ByteRate) {
		callback(rate.Mbps())
	})

	if err := st.server.UploadTestContext(ctx); err != nil {
		st.client.SetCallbackUpload(nil)
		return 0, fmt.Errorf("upload test failed: %w", err)
	}

	st.client.SetCallbackUpload(nil)
	return st.server.ULSpeed.Mbps(), nil
}
