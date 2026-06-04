package gui_wails

import (
	"context"
	"fmt"
	"log"

	"speedtest-tray/internal/speedtest_util"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) StartTest() {
	log.Println("Wails: StartTest called")

	updateCh := make(chan speedtest_util.Update)
	resultCh, err := a.tester.RunTest(context.Background(), updateCh)
	if err != nil {
		log.Printf("Wails: Failed to start test: %v\n", err)
		wailsRuntime.EventsEmit(a.ctx, "test_error", err.Error())
		return
	}

	go a.forwardTestEvents(updateCh, resultCh)
}

func (a *App) forwardTestEvents(updateCh <-chan speedtest_util.Update, resultCh <-chan speedtest_util.Result) {
	for update := range updateCh {
		wailsRuntime.EventsEmit(a.ctx, "test_update", serializeUpdate(update))
	}

	result := <-resultCh
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
