package gui_wails

import (
	"context"
	"log/slog"
	"time"

	"speedtest-tray/internal/config"
	"speedtest-tray/internal/speedtest_util"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type TestAdapter struct {
	ctx           context.Context
	tester        speedtest_util.TestOrchestrator
	emit          func(context.Context, string, ...interface{})
	resultTimeout time.Duration
}

func NewTestAdapter(ctx context.Context, tester speedtest_util.TestOrchestrator) *TestAdapter {
	return &TestAdapter{
		ctx:           ctx,
		tester:        tester,
		emit:          wailsRuntime.EventsEmit,
		resultTimeout: config.ResultTimeout,
	}
}

func (ta *TestAdapter) RunTest(ctx context.Context) (<-chan speedtest_util.Result, error) {
	updateCh := make(chan speedtest_util.Update, config.UpdateChannelSize)
	resultCh, err := speedtest_util.NewTestRunner(ta.tester).RunTest(ctx, updateCh)
	if err != nil {
		return nil, err
	}

	go ta.forwardUpdates(updateCh, resultCh)

	return resultCh, nil
}

func (ta *TestAdapter) forwardUpdates(updateCh <-chan speedtest_util.Update, resultCh <-chan speedtest_util.Result) {
	for update := range updateCh {
		slog.Debug(config.LogAdapterUpdate, "phase", update.Phase, "progress", update.Progress)
		event := serializeUpdate(update)
		ta.emit(ta.ctx, "test_update", event)
	}

	slog.Info(config.LogAdapterClosed)

	select {
	case result := <-resultCh:
		slog.Info(config.LogAdapterResult, "error", result.Error)
		if result.Error == nil {
			if err := speedtest_util.SaveToHistory(result.Server, result.Ping, result.Download, result.Upload); err != nil {
				slog.Error("Failed to save to history", "error", err)
			}
		}
		event := serializeResult(result)
		ta.emit(ta.ctx, "test_complete", event)
	case <-time.After(ta.resultTimeout):
		slog.Info(config.LogAdapterTimeout)
		ta.emit(ta.ctx, "test_complete", map[string]interface{}{"error": config.ErrTestTimeout})
	}
}

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

func formatNumber(value float64, precision int) string {
	return speedtest_util.FormatNumber(value, precision)
}

func failureStatus(update speedtest_util.Update) string {
	if update.Error != nil && update.Error.Error() == config.ErrNoInternet {
		return config.MsgNoInternet
	}
	return update.Phase.String()
}
