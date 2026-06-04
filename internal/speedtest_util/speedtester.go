package speedtest_util

import (
	"context"
	"fmt"
	"log"
	"time"

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

func (st *SpeedTester) RunTest(ctx context.Context, updateCh chan<- Update) (<-chan Result, error) {
	resultCh := make(chan Result, 1)
	log.Println("SpeedTester: RunTest requested")

	go func() {
		defer func() {
			log.Println("SpeedTester: Test goroutine exiting")
			close(updateCh)
			close(resultCh)
		}()

		// Helper to check if context was cancelled
		checkCancel := func() bool {
			if ctx.Err() != nil {
				log.Println("SpeedTester: Context cancelled, aborting test")
				st.fail(fmt.Errorf("Test stopped"), &Result{}, resultCh, updateCh)
				return true
			}
			return false
		}

		// Helper for interruptible sleep
		sleepWithContext := func(d time.Duration) bool {
			select {
			case <-ctx.Done():
				return checkCancel()
			case <-time.After(d):
				return false
			}
		}

		st.client.Reset()
		finalResult := Result{}

		log.Println("SpeedTester: Initializing...")
		updateCh <- Update{Phase: INITIALIZING, Progress: 0.05}
		if err := st.initialize(ctx, updateCh); err != nil {
			if ctx.Err() == nil { st.fail(err, &finalResult, resultCh, updateCh) }
			return
		}
		if checkCancel() { return }

		log.Println("SpeedTester: Selecting best server...")
		updateCh <- Update{Phase: SELECTING_SERVER, Progress: 0.10}
		if err := st.selectBestServer(ctx, updateCh, &finalResult); err != nil {
			if ctx.Err() == nil { st.fail(err, &finalResult, resultCh, updateCh) }
			return
		}
		if sleepWithContext(1 * time.Second) { return }

		log.Println("SpeedTester: Running ping test...")
		updateCh <- Update{Phase: PING_TEST, Progress: 0.20}
		if err := st.runPingTest(ctx, updateCh, &finalResult); err != nil {
			if ctx.Err() == nil { st.fail(err, &finalResult, resultCh, updateCh) }
			return
		}
		if sleepWithContext(1 * time.Second) { return }

		log.Println("SpeedTester: Starting download...")
		updateCh <- Update{Phase: STARTING_DOWNLOAD, Progress: 0.30, Ping: finalResult.Ping}
		if err := st.runDownloadTest(ctx, updateCh, &finalResult); err != nil {
			if ctx.Err() == nil { st.fail(err, &finalResult, resultCh, updateCh) }
			return
		}
		if sleepWithContext(1 * time.Second) { return }

		log.Println("SpeedTester: Starting upload...")
		updateCh <- Update{Phase: STARTING_UPLOAD, Progress: 0.70, Ping: finalResult.Ping, Download: finalResult.Download}
		if err := st.runUploadTest(ctx, updateCh, &finalResult); err != nil {
			if ctx.Err() == nil { st.fail(err, &finalResult, resultCh, updateCh) }
			return
		}

		log.Printf("SpeedTester: Test completed: Ping=%.2fms, DL=%.2fMbps, UL=%.2fMbps\n", finalResult.Ping, finalResult.Download, finalResult.Upload)
		updateCh <- Update{
			Phase:    COMPLETED,
			Progress: 1.0,
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
