package speedtest_util

import (
	"context"
	"fmt"
	"time"

	"github.com/showwin/speedtest-go/speedtest"
)

func (st *SpeedTester) runPingTest(ctx context.Context, updateCh chan<- Update, res *Result) error {
	updateCh <- Update{Phase: PING_TEST, Progress: 0.20}
	if err := st.server.PingTestContext(ctx, nil); err != nil {
		return fmt.Errorf("ping test failed: %w", err)
	}

	res.Ping = float64(st.server.Latency.Milliseconds())
	updateCh <- Update{
		Phase:    PING_TEST,
		Progress: 0.30,
		Ping:     res.Ping,
	}
	return nil
}

func (st *SpeedTester) runDownloadTest(ctx context.Context, updateCh chan<- Update, res *Result) error {
	updateCh <- Update{Phase: DOWNLOADING, Progress: 0.30, Ping: res.Ping}

	startTime := time.Now()
	testDuration := 10 * time.Second // Typical speedtest duration

	st.client.SetCallbackDownload(func(downRate speedtest.ByteRate) {
		elapsed := time.Since(startTime)
		// Calculate progress within the downloading phase (0.30 to 0.70)
		phaseProgress := float64(elapsed) / float64(testDuration)
		if phaseProgress > 1.0 {
			phaseProgress = 1.0
		}
		totalProgress := 0.30 + (phaseProgress * 0.40)

		updateCh <- Update{
			Phase:    DOWNLOADING,
			Progress: totalProgress,
			Ping:     res.Ping,
			Download: downRate.Mbps(),
		}
	})

	if err := st.server.DownloadTestContext(ctx); err != nil {
		return fmt.Errorf("download test failed: %w", err)
	}

	res.Download = st.server.DLSpeed.Mbps()
	return nil
}

func (st *SpeedTester) runUploadTest(ctx context.Context, updateCh chan<- Update, res *Result) error {
	updateCh <- Update{
		Phase:    UPLOADING,
		Progress: 0.70,
		Ping:     res.Ping,
		Download: res.Download,
	}

	startTime := time.Now()
	testDuration := 10 * time.Second // Typical speedtest duration

	st.client.SetCallbackUpload(func(upRate speedtest.ByteRate) {
		elapsed := time.Since(startTime)
		// Calculate progress within the uploading phase (0.70 to 0.95)
		phaseProgress := float64(elapsed) / float64(testDuration)
		if phaseProgress > 1.0 {
			phaseProgress = 1.0
		}
		totalProgress := 0.70 + (phaseProgress * 0.25)

		updateCh <- Update{
			Phase:    UPLOADING,
			Progress: totalProgress,
			Ping:     res.Ping,
			Download: res.Download,
			Upload:   upRate.Mbps(),
		}
	})

	if err := st.server.UploadTestContext(ctx); err != nil {
		return fmt.Errorf("upload test failed: %w", err)
	}

	res.Upload = st.server.ULSpeed.Mbps()
	return nil
}
