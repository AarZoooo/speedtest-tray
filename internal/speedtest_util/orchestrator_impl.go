package speedtest_util

import (
	"context"
	"fmt"
	"time"

	"speedtest-tray/internal/config"

	"github.com/showwin/speedtest-go/speedtest"
)

func (st *SpeedTester) GetUserInfo(ctx context.Context) error {
	user, err := st.client.FetchUserInfoContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch user info: %w", err)
	}
	st.user = user
	return nil
}

func (st *SpeedTester) FindServers(ctx context.Context) error {
	servers, err := st.client.FetchServers()
	if err != nil {
		return fmt.Errorf("failed to fetch servers: %w", err)
	}
	st.servers = servers
	return nil
}

func (st *SpeedTester) SelectBestServer(ctx context.Context) (*ServerInfo, error) {
	var targetIDs []int
	if st.TargetServerID != "" {
		var id int
		if _, err := fmt.Sscanf(st.TargetServerID, "%d", &id); err == nil {
			targetIDs = append(targetIDs, id)
		}
	}
	targets, err := st.servers.FindServer(targetIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to find target server: %w", err)
	}
	if len(targets) == 0 {
		return nil, fmt.Errorf("no server found with ID %s", st.TargetServerID)
	}
	st.server = targets[0]
	return &ServerInfo{
		Name:    st.server.Name,
		Country: st.server.Country,
	}, nil
}


func (st *SpeedTester) RunPing(ctx context.Context) (time.Duration, error) {
	if err := st.server.PingTestContext(ctx, nil); err != nil {
		return 0, fmt.Errorf("ping test failed: %w", err)
	}
	return st.server.Latency, nil
}

func (st *SpeedTester) RunDownload(ctx context.Context, callback func(float64)) (float64, error) {
	st.client.SetCaptureTime(config.EstimatedDurationDownload)
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

func (st *SpeedTester) RunUpload(ctx context.Context, callback func(float64)) (float64, error) {
	st.client.SetCaptureTime(config.EstimatedDurationUpload)
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
