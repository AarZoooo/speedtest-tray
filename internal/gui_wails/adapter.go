package gui_wails

import (
	"context"
	"log"
	"time"

	"speedtest-tray/internal/config"
	"speedtest-tray/internal/speedtest_util"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// TestAdapter bridges Wails events and speedtest logic
type TestAdapter struct {
	ctx           context.Context
	tester        speedtest_util.TestOrchestrator
	emit          func(context.Context, string, ...interface{})
	resultTimeout time.Duration
}

// NewTestAdapter creates a new adapter
func NewTestAdapter(ctx context.Context, tester speedtest_util.TestOrchestrator) *TestAdapter {
	return &TestAdapter{
		ctx:           ctx,
		tester:        tester,
		emit:          wailsRuntime.EventsEmit,
		resultTimeout: config.ResultTimeout,
	}
}

// RunTest orchestrates a speed test and emits Wails events
func (ta *TestAdapter) RunTest(ctx context.Context) (<-chan speedtest_util.Result, error) {
	updateCh := make(chan speedtest_util.Update, config.UpdateChannelSize)
	resultCh, err := speedtest_util.NewTestRunner(ta.tester).RunTest(ctx, updateCh)
	if err != nil {
		return nil, err
	}

	// Forward updates to Wails events
	go ta.forwardUpdates(updateCh, resultCh)

	return resultCh, nil
}

// forwardUpdates converts speedtest updates to Wails events
func (ta *TestAdapter) forwardUpdates(updateCh <-chan speedtest_util.Update, resultCh <-chan speedtest_util.Result) {
	for update := range updateCh {
		log.Printf("Adapter: Update - Phase=%s, Progress=%.2f\n", update.Phase, update.Progress)
		event := serializeUpdate(update)
		ta.emit(ta.ctx, "test_update", event)
	}

	log.Println("Adapter: Updates closed, waiting for result")

	select {
	case result := <-resultCh:
		log.Printf("Adapter: Result received - Error=%v\n", result.Error)
		event := serializeResult(result)
		ta.emit(ta.ctx, "test_complete", event)
	case <-time.After(ta.resultTimeout):
		log.Println("Adapter: Timeout waiting for result")
		ta.emit(ta.ctx, "test_complete", map[string]interface{}{"error": config.ErrTestTimeout})
	}
}

// serializeUpdate converts an Update to a Wails-compatible map
func serializeUpdate(update speedtest_util.Update) map[string]interface{} {
	return map[string]interface{}{
		"phase":    update.Phase,
		"status":   failureStatus(update),
		"progress": update.Progress,
		"ping":     formatNumber(update.Ping, 0),
		"download": formatNumber(update.Download, 1),
		"upload":   formatNumber(update.Upload, 1),
		"server":   update.Server,
	}
}

// serializeResult converts a Result to a Wails-compatible map
func serializeResult(result speedtest_util.Result) map[string]interface{} {
	if result.Error != nil {
		return map[string]interface{}{"error": result.Error.Error()}
	}

	return map[string]interface{}{
		"server":   result.Server,
		"ping":     formatNumber(result.Ping, 0),
		"download": formatNumber(result.Download, 1),
		"upload":   formatNumber(result.Upload, 1),
	}
}

// formatNumber formats a float with specified precision
func formatNumber(value float64, precision int) string {
	return speedtest_util.FormatNumber(value, precision)
}

func failureStatus(update speedtest_util.Update) string {
	if update.Error != nil && update.Error.Error() == config.ErrNoInternet {
		return config.MsgNoInternet
	}
	return update.Phase.String()
}
