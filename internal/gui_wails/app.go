package gui_wails

import (
	"context"

	"speedtest-tray/internal/speedtest_util"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App represents the Wails application binding
type App struct {
	ctx           context.Context
	cancel        context.CancelFunc
	windowVisible bool
	MacIcon       []byte
}

// NewApp creates a new App instance
func NewApp() *App {
	return &App{}
}

// Startup initializes the app on Wails startup
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.initMacStatusItem()
}

// Shutdown performs cleanup tasks when the application is closing
func (a *App) Shutdown(ctx context.Context) {
	a.cleanupMacStatusItem()
}

// ShowWindow displays the main window
func (a *App) ShowWindow() {
	if a.ctx == nil {
		return
	}

	a.positionWindow()
	wailsRuntime.WindowShow(a.ctx)
	wailsRuntime.WindowUnminimise(a.ctx)
	a.ApplyRoundedCorners()
	wailsRuntime.EventsEmit(a.ctx, "window_shown")
	a.windowVisible = true
}

// HideWindow hides the main window
func (a *App) HideWindow() {
	if a.ctx != nil {
		wailsRuntime.WindowHide(a.ctx)
		a.windowVisible = false
	}
}

// ToggleWindow toggles the main window visibility
func (a *App) ToggleWindow() {
	if a.windowVisible {
		a.HideWindow()
	} else {
		a.ShowWindow()
	}
}

// Quit quits the application
func (a *App) Quit() {
	if a.ctx != nil {
		wailsRuntime.Quit(a.ctx)
	}
}

// StartTest starts a speed test
func (a *App) StartTest() {
	a.StopTest()

	// Create a cancellable context for this specific test run
	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel

	tester := speedtest_util.New()
	adapter := NewTestAdapter(a.ctx, tester)
	_, err := adapter.RunTest(ctx)
	if err != nil {
		wailsRuntime.EventsEmit(a.ctx, "test_error", err.Error())
	}
}

// StopTest stops the running test
func (a *App) StopTest() {
	if a.cancel != nil {
		a.cancel()
		a.cancel = nil
	}
}
