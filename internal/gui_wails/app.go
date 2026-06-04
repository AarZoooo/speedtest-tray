package gui_wails

import (
	"context"

	"speedtest-tray/internal/config"
	"speedtest-tray/internal/speedtest_util"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App represents the Wails application binding
type App struct {
	ctx     context.Context
	adapter *TestAdapter
	tester  *speedtest_util.SpeedTester
	cancel  context.CancelFunc
}

// NewApp creates a new App instance
func NewApp(tester *speedtest_util.SpeedTester) *App {
	return &App{
		tester: tester,
	}
}

// Startup initializes the app on Wails startup
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.adapter = NewTestAdapter(ctx, a.tester)
}

// ShowWindow displays the main window
func (a *App) ShowWindow() {
	if a.ctx == nil {
		return
	}

	wailsRuntime.WindowSetSize(a.ctx, config.WindowWidth, config.WindowHeight)
	a.positionWindow()
	wailsRuntime.WindowShow(a.ctx)
	wailsRuntime.WindowUnminimise(a.ctx)
	a.ApplyRoundedCorners()
	wailsRuntime.EventsEmit(a.ctx, "window_shown")
}

// HideWindow hides the main window
func (a *App) HideWindow() {
	if a.ctx != nil {
		wailsRuntime.WindowHide(a.ctx)
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
	_, err := a.adapter.RunTest()
	if err != nil {
		wailsRuntime.EventsEmit(a.ctx, "test_error", err.Error())
	}
}

// StopTest stops the running test
func (a *App) StopTest() {
	if a.adapter != nil && a.cancel != nil {
		a.cancel()
	}
}
