package gui_wails

import (
	"context"
	"errors"
	"testing"
	"time"

	"speedtest-tray/internal/config"
	"speedtest-tray/internal/speedtest_util"
)

func TestSerializeUpdate(t *testing.T) {
	got := serializeUpdate(speedtest_util.Update{
		Phase:    speedtest_util.DOWNLOADING,
		Progress: 0.42,
		Ping:     23,
		Download: 91.54,
		Upload:   18.24,
		Server:   "Nearest (IN)",
	})

	assertMapValue(t, got, "phase", speedtest_util.DOWNLOADING)
	assertMapValue(t, got, "status", "Running download test...")
	assertMapValue(t, got, "progress", 0.42)
	assertMapValue(t, got, "ping", "23")
	assertMapValue(t, got, "download", "91.5")
	assertMapValue(t, got, "upload", "18.2")
	assertMapValue(t, got, "server", "Nearest (IN)")
}

func TestSerializeResultSuccess(t *testing.T) {
	got := serializeResult(speedtest_util.Result{
		Server:   "Nearest (IN)",
		Ping:     23,
		Download: 91.54,
		Upload:   18.24,
	})

	assertMapValue(t, got, "server", "Nearest (IN)")
	assertMapValue(t, got, "ping", "23")
	assertMapValue(t, got, "download", "91.5")
	assertMapValue(t, got, "upload", "18.2")
}

func TestSerializeResultError(t *testing.T) {
	got := serializeResult(speedtest_util.Result{Error: errors.New("network failed")})

	assertMapValue(t, got, "error", "network failed")
	if len(got) != 1 {
		t.Fatalf("serializeResult() = %+v, want only error", got)
	}
}

func TestForwardUpdatesEmitsUpdatesAndCompletion(t *testing.T) {
	recorder := &recordingEmitter{}
	adapter := NewTestAdapter(context.Background(), nil)
	adapter.emit = recorder.emit

	updateCh := make(chan speedtest_util.Update, 2)
	resultCh := make(chan speedtest_util.Result, 1)
	updateCh <- speedtest_util.Update{Phase: speedtest_util.INITIALIZING, Progress: config.ProgressInit}
	updateCh <- speedtest_util.Update{Phase: speedtest_util.COMPLETED, Progress: config.ProgressComplete}
	close(updateCh)
	resultCh <- speedtest_util.Result{Server: "Nearest (IN)", Ping: 23, Download: 91.5, Upload: 18.2}

	adapter.forwardUpdates(updateCh, resultCh)

	if len(recorder.events) != 3 {
		t.Fatalf("emitted %d events, want 3: %+v", len(recorder.events), recorder.events)
	}
	if recorder.events[0].name != "test_update" || recorder.events[1].name != "test_update" {
		t.Fatalf("first events = %+v, want test_update events", recorder.events[:2])
	}
	if recorder.events[2].name != "test_complete" {
		t.Fatalf("last event = %q, want test_complete", recorder.events[2].name)
	}
}

func TestForwardUpdatesEmitsTimeoutCompletion(t *testing.T) {
	recorder := &recordingEmitter{}
	adapter := NewTestAdapter(context.Background(), nil)
	adapter.emit = recorder.emit
	adapter.resultTimeout = time.Nanosecond

	updateCh := make(chan speedtest_util.Update)
	resultCh := make(chan speedtest_util.Result)
	close(updateCh)

	adapter.forwardUpdates(updateCh, resultCh)

	if len(recorder.events) != 1 {
		t.Fatalf("emitted %d events, want 1", len(recorder.events))
	}
	if recorder.events[0].name != "test_complete" {
		t.Fatalf("event name = %q, want test_complete", recorder.events[0].name)
	}
	payload, ok := recorder.events[0].args[0].(map[string]interface{})
	if !ok {
		t.Fatalf("payload type = %T, want map[string]interface{}", recorder.events[0].args[0])
	}
	assertMapValue(t, payload, "error", config.ErrTestTimeout)
}

type recordedEvent struct {
	name string
	args []interface{}
}

type recordingEmitter struct {
	events []recordedEvent
}

func (r *recordingEmitter) emit(_ context.Context, name string, args ...interface{}) {
	r.events = append(r.events, recordedEvent{name: name, args: args})
}

func assertMapValue(t *testing.T, got map[string]interface{}, key string, want interface{}) {
	t.Helper()
	if got[key] != want {
		t.Fatalf("%s = %#v, want %#v", key, got[key], want)
	}
}
