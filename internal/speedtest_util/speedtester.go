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

		log.Println("SpeedTester: Initializing...")
		updateCh <- Update{Phase: INITIALIZING, Progress: config.ProgressInit}
		if err := st.initialize(ctx, updateCh); err != nil {
			if ctx.Err() == nil { st.fail(err, &finalResult, resultCh, updateCh) }
			return
		}
		if st.checkContextCancelled(ctx, resultCh, updateCh) { return }

		log.Println("SpeedTester: Selecting best server...")
		updateCh <- Update{Phase: SELECTING_SERVER, Progress: config.ProgressSelectServer}
		if err := st.selectBestServer(ctx, updateCh, &finalResult); err != nil {
			if ctx.Err() == nil { st.fail(err, &finalResult, resultCh, updateCh) }
			return
		}
		if st.sleepWithInterrupt(ctx, 1*time.Second, resultCh, updateCh) { return }

		log.Println("SpeedTester: Running ping test...")
		updateCh <- Update{Phase: PING_TEST, Progress: config.ProgressPingStart}
		if err := st.runPingTest(ctx, updateCh, &finalResult); err != nil {
			if ctx.Err() == nil { st.fail(err, &finalResult, resultCh, updateCh) }
			return
		}
		if st.sleepWithInterrupt(ctx, 1*time.Second, resultCh, updateCh) { return }

		log.Println("SpeedTester: Starting download...")
		updateCh <- Update{Phase: STARTING_DOWNLOAD, Progress: config.ProgressDownStart, Ping: finalResult.Ping}
		if err := st.runDownloadTest(ctx, updateCh, &finalResult); err != nil {
			if ctx.Err() == nil { st.fail(err, &finalResult, resultCh, updateCh) }
			return
		}
		if st.sleepWithInterrupt(ctx, 1*time.Second, resultCh, updateCh) { return }

		log.Println("SpeedTester: Starting upload...")
		updateCh <- Update{Phase: STARTING_UPLOAD, Progress: config.ProgressUpStart, Ping: finalResult.Ping, Download: finalResult.Download}
		if err := st.runUploadTest(ctx, updateCh, &finalResult); err != nil {
			if ctx.Err() == nil { st.fail(err, &finalResult, resultCh, updateCh) }
			return
		}

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
