package speedtest_util

import (
	"context"
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
	log.Println("Starting speed test...")

	go func() {
		defer close(updateCh)
		defer close(resultCh)

		st.client.Reset()
		finalResult := Result{}

		log.Println("Initializing speedtest client...")
		if err := st.initialize(ctx, updateCh); err != nil {
			st.fail(err, &finalResult, resultCh, updateCh)
			return
		}
		time.Sleep(1 * time.Second)

		log.Println("Selecting best server...")
		if err := st.selectBestServer(ctx, updateCh, &finalResult); err != nil {
			st.fail(err, &finalResult, resultCh, updateCh)
			return
		}
		time.Sleep(1 * time.Second)

		log.Println("Running ping test...")
		if err := st.runPingTest(ctx, updateCh, &finalResult); err != nil {
			st.fail(err, &finalResult, resultCh, updateCh)
			return
		}
		updateCh <- Update{Phase: STARTING_DOWNLOAD, Progress: 0.30, Ping: finalResult.Ping}
		time.Sleep(1 * time.Second)

		log.Println("Running download test...")
		if err := st.runDownloadTest(ctx, updateCh, &finalResult); err != nil {
			st.fail(err, &finalResult, resultCh, updateCh)
			return
		}
		updateCh <- Update{Phase: STARTING_UPLOAD, Progress: 0.70, Ping: finalResult.Ping, Download: finalResult.Download}
		time.Sleep(1 * time.Second)

		log.Println("Running upload test...")
		if err := st.runUploadTest(ctx, updateCh, &finalResult); err != nil {
			st.fail(err, &finalResult, resultCh, updateCh)
			return
		}
		time.Sleep(1 * time.Second)

		log.Printf("Test completed: Ping=%.2fms, DL=%.2fMbps, UL=%.2fMbps\n", finalResult.Ping, finalResult.Download, finalResult.Upload)
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
