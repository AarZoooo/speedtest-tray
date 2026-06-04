package gui_wails

import (
	"context"
	"fmt"
	"log"
	"time"

	"speedtest-tray/internal/config"
	"speedtest-tray/internal/speedtest_util"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) StartTest() {
	log.Println("Wails: StartTest called")

	// Cancel any existing test
	if a.cancelFunc != nil {
		a.cancelFunc()
	}

	ctx, cancel := context.WithCancel(context.Background())
	a.cancelFunc = cancel

	updateCh := make(chan speedtest_util.Update)
	resultCh, err := a.tester.RunTest(ctx, updateCh)
	if err != nil {
		log.Printf("Wails: Failed to start test: %v\n", err)
		wailsRuntime.EventsEmit(a.ctx, "test_error", err.Error())
		return
	}

	log.Println("Wails: Test started successfully, forwarding events")
	go a.forwardTestEvents(updateCh, resultCh)
}

func (a *App) StopTest() {
	log.Println("Wails: StopTest called")
	if a.cancelFunc != nil {
		a.cancelFunc()
		a.cancelFunc = nil
	}
}

func (a *App) forwardTestEvents(updateCh <-chan speedtest_util.Update, resultCh <-chan speedtest_util.Result) {
	for update := range updateCh {
		log.Printf("Wails: Test update received: Phase=%s, Progress=%.2f\n", update.Phase, update.Progress)
		wailsRuntime.EventsEmit(a.ctx, "test_update", serializeUpdate(update))
	}

	log.Println("Wails: updateCh closed, waiting for final result")

	// Wait for the result with a timeout to avoid hanging if the test goroutine got stuck
	select {
	case result := <-resultCh:
		log.Printf("Wails: Final result received: Error=%v\n", result.Error)

		event := map[string]interface{}{
			"server":   result.Server,
			"ping":     formatNumber(result.Ping, 0),
			"download": formatNumber(result.Download, 1),
			"upload":   formatNumber(result.Upload, 1),
		}

		if result.Error != nil {
			event = map[string]interface{}{"error": result.Error.Error()}
		}

		wailsRuntime.EventsEmit(a.ctx, "test_complete", event)
		log.Println("Wails: test_complete event emitted")
	case <-time.After(config.ResultTimeout):
		log.Println("Wails: Timeout waiting for resultCh")
		wailsRuntime.EventsEmit(a.ctx, "test_complete", map[string]interface{}{"error": "Test stopped"})
	}
}

func serializeUpdate(update speedtest_util.Update) map[string]interface{} {
	return map[string]interface{}{
		"phase":    update.Phase,
		"status":   update.Phase.String(),
		"progress": update.Progress,
		"ping":     formatNumber(update.Ping, 0),
		"download": formatNumber(update.Download, 1),
		"upload":   formatNumber(update.Upload, 1),
		"server":   update.Server,
	}
}

func formatNumber(value float64, precision int) string {
	return fmt.Sprintf("%.*f", precision, value)
}
