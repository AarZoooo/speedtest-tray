package speedtest_util

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"speedtest-tray/internal/config"
)

type TestRunner struct {
	orchestrator  TestOrchestrator
	cancel        context.CancelFunc
	sleep         func(context.Context, time.Duration)
	checkInternet func(context.Context) error
	throttle      time.Duration
}

func NewTestRunner(orchestrator TestOrchestrator) *TestRunner {
	return &TestRunner{
		orchestrator:  orchestrator,
		sleep:         sleepOrCancel,
		checkInternet: CheckInternet,
		throttle:      100 * time.Millisecond,
	}
}

func (tr *TestRunner) RunTest(ctx context.Context, updateCh chan<- Update) (<-chan Result, error) {
	resultCh := make(chan Result, 1)
	ctx, tr.cancel = context.WithCancel(ctx)

	go func() {
		defer func() {
			slog.Info(config.LogHardwareStats, GetProcessStats().LogAttr()...)
			tr.cancel()
			close(updateCh)
			close(resultCh)
		}()

		slog.Info(config.LogHardwareStats, GetProcessStats().LogAttr()...)

		finalResult := Result{}

		tr.sendUpdate(ctx, updateCh, Update{Phase: INITIALIZING, Progress: config.ProgressInit})
		if err := tr.checkInternet(ctx); err != nil {
			tr.fail(err, &finalResult, resultCh, updateCh)
			return
		}
		if err := tr.orchestrator.GetUserInfo(ctx); err != nil {
			tr.fail(err, &finalResult, resultCh, updateCh)
			return
		}

		tr.sendUpdate(ctx, updateCh, Update{Phase: GETTING_INFO, Progress: config.ProgressGetInfo})
		if ctx.Err() != nil {
			tr.fail(fmt.Errorf(config.ErrTestStopped), &finalResult, resultCh, updateCh)
			return
		}

		tr.sendUpdate(ctx, updateCh, Update{Phase: FINDING_SERVERS, Progress: config.ProgressFindServers})
		if err := tr.orchestrator.FindServers(ctx); err != nil {
			tr.fail(err, &finalResult, resultCh, updateCh)
			return
		}

		tr.sendUpdate(ctx, updateCh, Update{Phase: SELECTING_SERVER, Progress: config.ProgressSelectServer})
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
		tr.sendUpdate(ctx, updateCh, Update{Phase: SERVER_SELECTED, Progress: config.ProgressServerSelect, Server: finalResult.Server})
		tr.sleep(ctx, config.PhaseSleepDuration)
		if ctx.Err() != nil {
			tr.fail(fmt.Errorf(config.ErrTestStopped), &finalResult, resultCh, updateCh)
			return
		}

		tr.sendUpdate(ctx, updateCh, Update{Phase: PING_TEST, Progress: config.ProgressPingStart})
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
		tr.sendUpdate(ctx, updateCh, Update{Phase: PING_TEST, Progress: config.ProgressPingEnd, Ping: finalResult.Ping})
		tr.sendUpdate(ctx, updateCh, Update{Phase: STARTING_DOWNLOAD, Progress: config.ProgressDownStart, Ping: finalResult.Ping})
		tr.sleep(ctx, config.PhaseSleepDuration)
		if ctx.Err() != nil {
			tr.fail(fmt.Errorf(config.ErrTestStopped), &finalResult, resultCh, updateCh)
			return
		}
		dlStart := time.Now()
		dlResCh := make(chan struct {
			val float64
			err error
		}, 1)
		lastDlUpdate := time.Now().Add(-tr.throttle)

		go func() {
			val, err := tr.orchestrator.RunDownload(ctx, func(mbps float64) {
				if tr.throttle > 0 && time.Since(lastDlUpdate) < tr.throttle {
					return
				}
				lastDlUpdate = time.Now()

				elapsed := time.Since(dlStart).Seconds()
				duration := config.EstimatedDurationDownload.Seconds()
				progress := CalculatePhaseProgress(elapsed, duration)
				totalProgress := MapPhaseProgressToTotal(config.ProgressDownStart, config.ProgressDownEnd, progress)
				tr.sendUpdate(ctx, updateCh, Update{
					Phase:    DOWNLOADING,
					Progress: totalProgress,
					Ping:     finalResult.Ping,
					Download: mbps,
				})
			})
			dlResCh <- struct {
				val float64
				err error
			}{val, err}
		}()

		select {
		case res := <-dlResCh:
			if res.err != nil {
				tr.fail(res.err, &finalResult, resultCh, updateCh)
				return
			}
			finalResult.Download = res.val
		case <-ctx.Done():
			tr.fail(fmt.Errorf(config.ErrTestStopped), &finalResult, resultCh, updateCh)
			return
		}

		tr.sendUpdate(ctx, updateCh, Update{Phase: STARTING_UPLOAD, Progress: config.ProgressUpStart, Ping: finalResult.Ping, Download: finalResult.Download})
		tr.sleep(ctx, config.PhaseSleepDuration)
		if ctx.Err() != nil {
			tr.fail(fmt.Errorf(config.ErrTestStopped), &finalResult, resultCh, updateCh)
			return
		}
		ulStart := time.Now()
		ulResCh := make(chan struct {
			val float64
			err error
		}, 1)
		lastUlUpdate := time.Now().Add(-tr.throttle)

		go func() {
			val, err := tr.orchestrator.RunUpload(ctx, func(mbps float64) {
				if tr.throttle > 0 && time.Since(lastUlUpdate) < tr.throttle {
					return
				}
				lastUlUpdate = time.Now()

				elapsed := time.Since(ulStart).Seconds()
				duration := config.EstimatedDurationUpload.Seconds()
				progress := CalculatePhaseProgress(elapsed, duration)
				totalProgress := MapPhaseProgressToTotal(config.ProgressUpStart, config.ProgressUpEnd, progress)
				tr.sendUpdate(ctx, updateCh, Update{
					Phase:    UPLOADING,
					Progress: totalProgress,
					Ping:     finalResult.Ping,
					Download: finalResult.Download,
					Upload:   mbps,
				})
			})
			ulResCh <- struct {
				val float64
				err error
			}{val, err}
		}()

		select {
		case res := <-ulResCh:
			if res.err != nil {
				tr.fail(res.err, &finalResult, resultCh, updateCh)
				return
			}
			finalResult.Upload = res.val
		case <-ctx.Done():
			tr.fail(fmt.Errorf(config.ErrTestStopped), &finalResult, resultCh, updateCh)
			return
		}

		tr.sendUpdate(ctx, updateCh, Update{
			Phase:    COMPLETED,
			Progress: config.ProgressComplete,
			Ping:     finalResult.Ping,
			Download: finalResult.Download,
			Upload:   finalResult.Upload,
		})
		finalResult.Phase = COMPLETED
		resultCh <- finalResult
	}()

	return resultCh, nil
}

func (tr *TestRunner) Cancel() {
	if tr.cancel != nil {
		tr.cancel()
	}
}

func (tr *TestRunner) sendUpdate(ctx context.Context, updateCh chan<- Update, update Update) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	defer func() {
		_ = recover()
	}()
	updateCh <- update
}

func (tr *TestRunner) fail(err error, res *Result, resCh chan<- Result, updateCh chan<- Update) {
	res.Error = err
	phase := FAILED
	if err.Error() == config.ErrTestStopped {
		phase = config.PhaseStopped
	}
	res.Phase = phase
	tr.sendUpdate(context.Background(), updateCh, Update{Phase: phase, Error: err})
	resCh <- *res
}

func sleepOrCancel(ctx context.Context, d time.Duration) {
	select {
	case <-ctx.Done():
		return
	case <-time.After(d):
		return
	}
}
