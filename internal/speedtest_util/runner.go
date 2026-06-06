package speedtest_util

import (
	"context"
	"fmt"
	"log"
	"time"

	"speedtest-tray/internal/config"
)

// TestRunner wraps TestOrchestrator and handles test orchestration with proper lifecycle
type TestRunner struct {
	orchestrator TestOrchestrator
	cancel       context.CancelFunc
	sleep        func(context.Context, time.Duration)
}

// NewTestRunner creates a new test runner
func NewTestRunner(orchestrator TestOrchestrator) *TestRunner {
	return &TestRunner{
		orchestrator: orchestrator,
		sleep:        sleepOrCancel,
	}
}

// RunTest orchestrates the full speed test workflow
func (tr *TestRunner) RunTest(ctx context.Context, updateCh chan<- Update) (<-chan Result, error) {
	resultCh := make(chan Result, 1)
	log.Println("TestRunner: RunTest requested")

	ctx, tr.cancel = context.WithCancel(ctx)

	go func() {
		defer func() {
			log.Println("TestRunner: Test goroutine exiting")
			close(updateCh)
			close(resultCh)
		}()

		finalResult := Result{}

		// Phase: Initialize
		log.Println("TestRunner: Initializing...")
		updateCh <- Update{Phase: INITIALIZING, Progress: config.ProgressInit}
		if err := tr.orchestrator.GetUserInfo(ctx); err != nil {
			tr.fail(err, &finalResult, resultCh, updateCh)
			return
		}
		updateCh <- Update{Phase: GETTING_INFO, Progress: config.ProgressGetInfo}
		if ctx.Err() != nil {
			tr.fail(fmt.Errorf(config.ErrTestStopped), &finalResult, resultCh, updateCh)
			return
		}

		// Phase: Find and select server
		log.Println("TestRunner: Locating servers...")
		updateCh <- Update{Phase: FINDING_SERVERS, Progress: config.ProgressFindServers}
		if err := tr.orchestrator.FindServers(ctx); err != nil {
			tr.fail(err, &finalResult, resultCh, updateCh)
			return
		}

		log.Println("TestRunner: Selecting best server...")
		updateCh <- Update{Phase: SELECTING_SERVER, Progress: config.ProgressSelectServer}
		serverInfo, err := tr.orchestrator.SelectBestServer(ctx)
		if err != nil {
			tr.fail(err, &finalResult, resultCh, updateCh)
			return
		}
		if ctx.Err() != nil {
			tr.fail(fmt.Errorf(config.ErrTestStopped), &finalResult, resultCh, updateCh)
			return
		}

		finalResult.Server = fmt.Sprintf("%s (%s)", serverInfo.Name, serverInfo.Country)
		log.Printf("Selected server: %s\n", finalResult.Server)
		updateCh <- Update{Phase: SERVER_SELECTED, Progress: config.ProgressServerSelect, Server: finalResult.Server}
		tr.sleep(ctx, 1*time.Second)
		if ctx.Err() != nil {
			tr.fail(fmt.Errorf(config.ErrTestStopped), &finalResult, resultCh, updateCh)
			return
		}

		// Phase: Ping test
		log.Println("TestRunner: Running ping test...")
		updateCh <- Update{Phase: PING_TEST, Progress: config.ProgressPingStart}
		latency, err := tr.orchestrator.RunPing(ctx)
		if err != nil {
			tr.fail(err, &finalResult, resultCh, updateCh)
			return
		}
		if ctx.Err() != nil {
			tr.fail(fmt.Errorf(config.ErrTestStopped), &finalResult, resultCh, updateCh)
			return
		}

		finalResult.Ping = float64(latency.Milliseconds())
		updateCh <- Update{Phase: PING_TEST, Progress: config.ProgressPingEnd, Ping: finalResult.Ping}
		tr.sleep(ctx, 1*time.Second)
		if ctx.Err() != nil {
			tr.fail(fmt.Errorf(config.ErrTestStopped), &finalResult, resultCh, updateCh)
			return
		}

		// Phase: Download test
		log.Println("TestRunner: Starting download test...")
		updateCh <- Update{Phase: STARTING_DOWNLOAD, Progress: config.ProgressDownStart, Ping: finalResult.Ping}
		dlStart := time.Now()
		downloadSpeed, err := tr.orchestrator.RunDownload(ctx, func(mbps float64) {
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
			tr.fail(err, &finalResult, resultCh, updateCh)
			return
		}
		if ctx.Err() != nil {
			tr.fail(fmt.Errorf(config.ErrTestStopped), &finalResult, resultCh, updateCh)
			return
		}

		finalResult.Download = downloadSpeed
		tr.sleep(ctx, 1*time.Second)
		if ctx.Err() != nil {
			tr.fail(fmt.Errorf(config.ErrTestStopped), &finalResult, resultCh, updateCh)
			return
		}

		// Phase: Upload test
		log.Println("TestRunner: Starting upload test...")
		updateCh <- Update{Phase: STARTING_UPLOAD, Progress: config.ProgressUpStart, Ping: finalResult.Ping, Download: finalResult.Download}
		ulStart := time.Now()
		uploadSpeed, err := tr.orchestrator.RunUpload(ctx, func(mbps float64) {
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
			tr.fail(err, &finalResult, resultCh, updateCh)
			return
		}
		if ctx.Err() != nil {
			tr.fail(fmt.Errorf(config.ErrTestStopped), &finalResult, resultCh, updateCh)
			return
		}

		finalResult.Upload = uploadSpeed

		// Complete
		log.Printf("TestRunner: Test completed: Ping=%.2fms, DL=%.2fMbps, UL=%.2fMbps\n", finalResult.Ping, finalResult.Download, finalResult.Upload)
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

// Cancel stops the running test
func (tr *TestRunner) Cancel() {
	if tr.cancel != nil {
		tr.cancel()
	}
}

// fail sends failure result
func (tr *TestRunner) fail(err error, res *Result, resCh chan<- Result, updateCh chan<- Update) {
	log.Printf("TestRunner: Test failed: %v\n", err)
	res.Error = err

	phase := FAILED
	if err.Error() == config.ErrTestStopped {
		phase = config.PhaseStopped
	}

	res.Phase = phase
	updateCh <- Update{Phase: phase, Error: err}
	resCh <- *res
}

// sleepOrCancel sleeps but returns early if context is cancelled
func sleepOrCancel(ctx context.Context, d time.Duration) {
	select {
	case <-ctx.Done():
		return
	case <-time.After(d):
		return
	}
}
