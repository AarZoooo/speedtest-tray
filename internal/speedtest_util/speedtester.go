package speedtest_util

import (
	"context"
	"fmt"
	"log"
	"time"

	"speedtest-tray/internal/config"

	"github.com/showwin/speedtest-go/speedtest"
)

type SpeedTester struct {
	client  *speedtest.Speedtest
	user    *speedtest.User
	servers speedtest.Servers
	server  *speedtest.Server
}

func New() *SpeedTester {
	return &SpeedTester{
		client: speedtest.New(),
	}
}

// checkContextCancelled checks if context is cancelled and handles cleanup
func (st *SpeedTester) checkContextCancelled(ctx context.Context, resultCh chan<- Result, updateCh chan<- Update) bool {
	if ctx.Err() != nil {
		log.Println("SpeedTester: Context cancelled, aborting test")
		st.fail(fmt.Errorf("Test stopped"), &Result{}, resultCh, updateCh)
		return true
	}
	return false
}

// sleepWithInterrupt sleeps for duration but can be interrupted by context cancellation
func (st *SpeedTester) sleepWithInterrupt(ctx context.Context, d time.Duration, resultCh chan<- Result, updateCh chan<- Update) bool {
	select {
	case <-ctx.Done():
		return st.checkContextCancelled(ctx, resultCh, updateCh)
	case <-time.After(d):
		return false
	}
}

func (st *SpeedTester) RunTest(ctx context.Context, updateCh chan<- Update) (<-chan Result, error) {
	resultCh := make(chan Result, 1)
	log.Println("SpeedTester: RunTest requested")

	go func() {
		defer func() {
			log.Println("SpeedTester: Test goroutine exiting")
			close(updateCh)
			close(resultCh)
		}()

		st.client.Reset()
		finalResult := Result{}

		// Phase: Initialize
		log.Println("SpeedTester: Initializing...")
		updateCh <- Update{Phase: INITIALIZING, Progress: config.ProgressInit}
		if err := st.GetUserInfo(ctx); err != nil {
			if ctx.Err() == nil { st.fail(err, &finalResult, resultCh, updateCh) }
			return
		}
		updateCh <- Update{Phase: GETTING_INFO, Progress: config.ProgressGetInfo}
		if st.checkContextCancelled(ctx, resultCh, updateCh) { return }

		// Phase: Find and select server
		log.Println("SpeedTester: Locating servers...")
		updateCh <- Update{Phase: FINDING_SERVERS, Progress: config.ProgressFindServers}
		if err := st.FindServers(ctx); err != nil {
			if ctx.Err() == nil { st.fail(err, &finalResult, resultCh, updateCh) }
			return
		}

		log.Println("SpeedTester: Selecting best server...")
		updateCh <- Update{Phase: SELECTING_SERVER, Progress: config.ProgressSelectServer}
		serverInfo, err := st.SelectBestServer(ctx)
		if err != nil {
			if ctx.Err() == nil { st.fail(err, &finalResult, resultCh, updateCh) }
			return
		}
		finalResult.Server = fmt.Sprintf("%s (%s)", serverInfo.Name, serverInfo.Country)
		log.Printf("Selected server: %s\n", finalResult.Server)
		updateCh <- Update{Phase: SERVER_SELECTED, Progress: config.ProgressServerSelect, Server: finalResult.Server}
		if st.sleepWithInterrupt(ctx, 1*time.Second, resultCh, updateCh) { return }

		// Phase: Ping test
		log.Println("SpeedTester: Running ping test...")
		updateCh <- Update{Phase: PING_TEST, Progress: config.ProgressPingStart}
		latency, err := st.RunPing(ctx)
		if err != nil {
			if ctx.Err() == nil { st.fail(err, &finalResult, resultCh, updateCh) }
			return
		}
		finalResult.Ping = float64(latency.Milliseconds())
		updateCh <- Update{Phase: PING_TEST, Progress: config.ProgressPingEnd, Ping: finalResult.Ping}
		if st.sleepWithInterrupt(ctx, 1*time.Second, resultCh, updateCh) { return }

		// Phase: Download test
		log.Println("SpeedTester: Starting download test...")
		updateCh <- Update{Phase: STARTING_DOWNLOAD, Progress: config.ProgressDownStart, Ping: finalResult.Ping}
		dlStart := time.Now()
		downloadSpeed, err := st.RunDownload(ctx, func(mbps float64) {
			elapsed := time.Since(dlStart).Seconds()
			duration := config.TestDurationDownload.Seconds()
			progress := CalculatePhaseProgress(elapsed, duration)
			totalProgress := MapPhaseProgressToTotal(config.ProgressDownStart, config.ProgressDownEnd, progress)
			updateCh <- Update{
				Phase:    DOWNLOADING,
				Progress: totalProgress,
				Ping:     finalResult.Ping,
				Download: mbps,
			}
		})
		if err != nil {
			if ctx.Err() == nil { st.fail(err, &finalResult, resultCh, updateCh) }
			return
		}
		finalResult.Download = downloadSpeed
		if st.sleepWithInterrupt(ctx, 1*time.Second, resultCh, updateCh) { return }

		// Phase: Upload test
		log.Println("SpeedTester: Starting upload test...")
		updateCh <- Update{Phase: STARTING_UPLOAD, Progress: config.ProgressUpStart, Ping: finalResult.Ping, Download: finalResult.Download}
		ulStart := time.Now()
		uploadSpeed, err := st.RunUpload(ctx, func(mbps float64) {
			elapsed := time.Since(ulStart).Seconds()
			duration := config.TestDurationUpload.Seconds()
			progress := CalculatePhaseProgress(elapsed, duration)
			totalProgress := MapPhaseProgressToTotal(config.ProgressUpStart, config.ProgressUpEnd, progress)
			updateCh <- Update{
				Phase:    UPLOADING,
				Progress: totalProgress,
				Ping:     finalResult.Ping,
				Download: finalResult.Download,
				Upload:   mbps,
			}
		})
		if err != nil {
			if ctx.Err() == nil { st.fail(err, &finalResult, resultCh, updateCh) }
			return
		}
		finalResult.Upload = uploadSpeed

		log.Printf("SpeedTester: Test completed: Ping=%.2fms, DL=%.2fMbps, UL=%.2fMbps\n", finalResult.Ping, finalResult.Download, finalResult.Upload)
		updateCh <- Update{
			Phase:    COMPLETED,
			Progress: config.ProgressComplete,
			Ping:     finalResult.Ping,
			Download: finalResult.Download,
			Upload:   finalResult.Upload,
		}
		finalResult.Phase = COMPLETED
		resultCh <- finalResult
	}()

	return resultCh, nil
}

func (st *SpeedTester) fail(err error, res *Result, resCh chan<- Result, updateCh chan<- Update) {
	log.Printf("Test failed: %v\n", err)
	res.Error = err
	res.Phase = FAILED
	updateCh <- Update{Phase: FAILED, Error: err}
	resCh <- *res
}
