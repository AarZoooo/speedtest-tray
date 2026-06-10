package speedtest_util

import (
	"context"
	"errors"
	"testing"
	"time"

	"speedtest-tray/internal/config"
)

func TestRunTestSuccess(t *testing.T) {
	mock := &MockOrchestrator{
		ServerResult:    &ServerInfo{Name: "Nearest", Country: "IN"},
		PingResult:      23 * time.Millisecond,
		DownloadResult:  91.5,
		UploadResult:    18.2,
		DownloadSamples: []float64{12, 48, 91.5},
		UploadSamples:   []float64{5, 12, 18.2},
	}

	result, updates := runTestWithMock(t, mock)

	if result.Error != nil {
		t.Fatalf("RunTest() result error = %v", result.Error)
	}
	if result.Phase != COMPLETED {
		t.Fatalf("result phase = %s, want %s", result.Phase, COMPLETED)
	}
	if result.Server != "Nearest (IN)" {
		t.Fatalf("result server = %q", result.Server)
	}
	if result.Ping != 23 || result.Download != 91.5 || result.Upload != 18.2 {
		t.Fatalf("unexpected result values: %+v", result)
	}

	mock.VerifyCalls(t, 1, 1, 1, 1, 1, 1)
	assertPhaseOrder(t, updates, []Phase{
		INITIALIZING,
		GETTING_INFO,
		FINDING_SERVERS,
		SELECTING_SERVER,
		SERVER_SELECTED,
		PING_TEST,
		PING_TEST,
		STARTING_DOWNLOAD,
		DOWNLOADING,
		DOWNLOADING,
		DOWNLOADING,
		STARTING_UPLOAD,
		UPLOADING,
		UPLOADING,
		UPLOADING,
		COMPLETED,
	})
}

func TestRunTestProgressMapping(t *testing.T) {
	mock := &MockOrchestrator{
		DownloadResult:  80,
		UploadResult:    20,
		DownloadSamples: []float64{10, 40, 80},
		UploadSamples:   []float64{5, 10, 20},
	}

	_, updates := runTestWithMock(t, mock)

	var previous float64
	for _, update := range updates {
		if update.Progress < previous {
			t.Fatalf("progress decreased from %v to %v at phase %s", previous, update.Progress, update.Phase)
		}
		previous = update.Progress

		switch update.Phase {
		case DOWNLOADING:
			if update.Progress < config.ProgressDownStart || update.Progress > config.ProgressDownEnd {
				t.Fatalf("download progress %v outside [%v, %v]", update.Progress, config.ProgressDownStart, config.ProgressDownEnd)
			}
		case UPLOADING:
			if update.Progress < config.ProgressUpStart || update.Progress > config.ProgressUpEnd {
				t.Fatalf("upload progress %v outside [%v, %v]", update.Progress, config.ProgressUpStart, config.ProgressUpEnd)
			}
		}
	}
}

func TestRunTestCanceledBeforeStart(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	mock := &MockOrchestrator{}

	result, updates := runTestWithContextAndMock(t, ctx, mock)

	if result.Error == nil || result.Error.Error() != config.ErrTestStopped {
		t.Fatalf("result error = %v, want %q", result.Error, config.ErrTestStopped)
	}
	if result.Phase != config.PhaseStopped {
		t.Fatalf("result phase = %s, want %s", result.Phase, config.PhaseStopped)
	}
	if last := updates[len(updates)-1]; last.Phase != config.PhaseStopped {
		t.Fatalf("last update phase = %s, want %s", last.Phase, config.PhaseStopped)
	}
	mock.VerifyCalls(t, 1, 0, 0, 0, 0, 0)
}

func TestRunTestCanceledDuringServerSelection(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	mock := &MockOrchestrator{
		OnSelectBestServer: func(context.Context) {
			cancel()
		},
	}

	result, _ := runTestWithContextAndMock(t, ctx, mock)

	if result.Error == nil || result.Error.Error() != config.ErrTestStopped {
		t.Fatalf("result error = %v, want %q", result.Error, config.ErrTestStopped)
	}
	mock.VerifyCalls(t, 1, 1, 1, 0, 0, 0)
}

func TestRunTestErrorsStopDownstreamPhases(t *testing.T) {
	tests := []struct {
		name        string
		configure   func(*MockOrchestrator, error)
		wantCalls   []int
		wantMessage string
	}{
		{
			name: "get_user_info",
			configure: func(m *MockOrchestrator, err error) {
				m.GetUserInfoErr = err
			},
			wantCalls: []int{1, 0, 0, 0, 0, 0},
		},
		{
			name: "find_servers",
			configure: func(m *MockOrchestrator, err error) {
				m.FindServersErr = err
			},
			wantCalls: []int{1, 1, 0, 0, 0, 0},
		},
		{
			name: "select_best_server",
			configure: func(m *MockOrchestrator, err error) {
				m.SelectBestServerErr = err
			},
			wantCalls: []int{1, 1, 1, 0, 0, 0},
		},
		{
			name: "run_ping",
			configure: func(m *MockOrchestrator, err error) {
				m.RunPingErr = err
			},
			wantCalls: []int{1, 1, 1, 1, 0, 0},
		},
		{
			name: "run_download",
			configure: func(m *MockOrchestrator, err error) {
				m.RunDownloadErr = err
			},
			wantCalls: []int{1, 1, 1, 1, 1, 0},
		},
		{
			name: "run_upload",
			configure: func(m *MockOrchestrator, err error) {
				m.RunUploadErr = err
			},
			wantCalls: []int{1, 1, 1, 1, 1, 1},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wantErr := errors.New(test.name)
			mock := &MockOrchestrator{}
			test.configure(mock, wantErr)

			result, updates := runTestWithMock(t, mock)

			if !errors.Is(result.Error, wantErr) {
				t.Fatalf("result error = %v, want %v", result.Error, wantErr)
			}
			if result.Phase != FAILED {
				t.Fatalf("result phase = %s, want %s", result.Phase, FAILED)
			}
			if updates[len(updates)-1].Phase != FAILED {
				t.Fatalf("last update phase = %s, want %s", updates[len(updates)-1].Phase, FAILED)
			}
			mock.VerifyCalls(t, test.wantCalls[0], test.wantCalls[1], test.wantCalls[2], test.wantCalls[3], test.wantCalls[4], test.wantCalls[5])
		})
	}
}

func TestRunTestNoInternet(t *testing.T) {
	mock := &MockOrchestrator{}
	updateCh := make(chan Update, 8)
	runner := NewTestRunner(mock)
	runner.sleep = func(context.Context, time.Duration) {}
	runner.checkInternet = func(context.Context) error {
		return errors.New(config.ErrNoInternet)
	}

	resultCh, err := runner.RunTest(context.Background(), updateCh)
	if err != nil {
		t.Fatalf("RunTest() error = %v", err)
	}

	result := <-resultCh
	if result.Error == nil || result.Error.Error() != config.ErrNoInternet {
		t.Fatalf("result error = %v, want %q", result.Error, config.ErrNoInternet)
	}
	if result.Phase != FAILED {
		t.Fatalf("result phase = %s, want %s", result.Phase, FAILED)
	}

	var updates []Update
	for update := range updateCh {
		updates = append(updates, update)
	}

	if len(updates) != 2 {
		t.Fatalf("got %d updates, want 2: %+v", len(updates), updates)
	}
	if updates[0].Phase != INITIALIZING {
		t.Fatalf("first phase = %s, want %s", updates[0].Phase, INITIALIZING)
	}
	if updates[1].Phase != FAILED {
		t.Fatalf("last phase = %s, want %s", updates[1].Phase, FAILED)
	}
	mock.VerifyCalls(t, 0, 0, 0, 0, 0, 0)
}

func TestRunTestCallbackAfterClose(t *testing.T) {
	var storedCallback func(float64)
	mock := &MockOrchestrator{
		OnRunDownload: func(ctx context.Context, cb func(float64)) {
			storedCallback = cb
		},
		RunDownloadErr: errors.New("abort"),
	}

	updateCh := make(chan Update, 10)
	runner := NewTestRunner(mock)
	runner.sleep = func(context.Context, time.Duration) {}
	runner.checkInternet = func(context.Context) error { return nil }

	resultCh, err := runner.RunTest(context.Background(), updateCh)
	if err != nil {
		t.Fatalf("RunTest() error = %v", err)
	}
	<-resultCh

	time.Sleep(50 * time.Millisecond)

	if storedCallback != nil {
		storedCallback(100.0)
	} else {
		t.Fatal("callback not stored")
	}
}

func runTestWithMock(t *testing.T, mock *MockOrchestrator) (Result, []Update) {
	t.Helper()
	return runTestWithContextAndMock(t, context.Background(), mock)
}

func runTestWithContextAndMock(t *testing.T, ctx context.Context, mock *MockOrchestrator) (Result, []Update) {
	t.Helper()
	updateCh := make(chan Update, 64)
	runner := NewTestRunner(mock)
	runner.sleep = func(context.Context, time.Duration) {}
	runner.checkInternet = func(context.Context) error { return nil }

	resultCh, err := runner.RunTest(ctx, updateCh)
	if err != nil {
		t.Fatalf("RunTest() error = %v", err)
	}

	result, ok := <-resultCh
	if !ok {
		t.Fatal("result channel closed without result")
	}

	var updates []Update
	for update := range updateCh {
		updates = append(updates, update)
	}

	if _, ok := <-resultCh; ok {
		t.Fatal("result channel should be closed after final result")
	}
	if len(updates) == 0 {
		t.Fatal("expected at least one update")
	}
	return result, updates
}

func assertPhaseOrder(t *testing.T, updates []Update, want []Phase) {
	t.Helper()
	if len(updates) != len(want) {
		t.Fatalf("got %d updates, want %d: %+v", len(updates), len(want), updates)
	}
	for i := range want {
		if updates[i].Phase != want[i] {
			t.Fatalf("update %d phase = %s, want %s", i, updates[i].Phase, want[i])
		}
	}
}
