package speedtest_util

import (
	"context"
	"fmt"
	"time"

	"speedtest-tray/internal/config"

	"github.com/showwin/speedtest-go/speedtest"
)

func (st *SpeedTester) runPingTest(ctx context.Context, updateCh chan<- Update, res *Result) error {
	updateCh <- Update{Phase: PING_TEST, Progress: config.ProgressPingStart}
	if err := st.server.PingTestContext(ctx, nil); err != nil {
		return fmt.Errorf("ping test failed: %w", err)
	}

	res.Ping = float64(st.server.Latency.Milliseconds())
	updateCh <- Update{
		Phase:    PING_TEST,
		Progress: config.ProgressPingEnd,
		Ping:     res.Ping,
	}
	return nil
}

func (st *SpeedTester) runDownloadTest(ctx context.Context, updateCh chan<- Update, res *Result) error {
	updateCh <- Update{Phase: DOWNLOADING, Progress: config.ProgressDownStart, Ping: res.Ping}

	startTime := time.Now()

	st.client.SetCallbackDownload(func(downRate speedtest.ByteRate) {
		elapsed := time.Since(startTime).Seconds()
		duration := config.TestDurationDownload.Seconds()
		phaseProgress := CalculatePhaseProgress(elapsed, duration)
		totalProgress := MapPhaseProgressToTotal(config.ProgressDownStart, config.ProgressDownEnd, phaseProgress)

		updateCh <- Update{
			Phase:    DOWNLOADING,
			Progress: totalProgress,
			Ping:     res.Ping,
			Download: downRate.Mbps(),
		}
	})

	if err := st.server.DownloadTestContext(ctx); err != nil {
		st.client.SetCallbackDownload(nil)
		return fmt.Errorf("download test failed: %w", err)
	}

	st.client.SetCallbackDownload(nil)
	res.Download = st.server.DLSpeed.Mbps()
	return nil
}

func (st *SpeedTester) runUploadTest(ctx context.Context, updateCh chan<- Update, res *Result) error {
	updateCh <- Update{
		Phase:    UPLOADING,
		Progress: config.ProgressUpStart,
		Ping:     res.Ping,
		Download: res.Download,
	}

	startTime := time.Now()

	st.client.SetCallbackUpload(func(upRate speedtest.ByteRate) {
		elapsed := time.Since(startTime).Seconds()
		duration := config.TestDurationUpload.Seconds()
		phaseProgress := CalculatePhaseProgress(elapsed, duration)
		totalProgress := MapPhaseProgressToTotal(config.ProgressUpStart, config.ProgressUpEnd, phaseProgress)

		updateCh <- Update{
			Phase:    UPLOADING,
			Progress: totalProgress,
			Ping:     res.Ping,
			Download: res.Download,
			Upload:   upRate.Mbps(),
		}
	})

	if err := st.server.UploadTestContext(ctx); err != nil {
		st.client.SetCallbackUpload(nil)
		return fmt.Errorf("upload test failed: %w", err)
	}

	st.client.SetCallbackUpload(nil)
	res.Upload = st.server.ULSpeed.Mbps()
	return nil
}
