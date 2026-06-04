package gui_wails

import (
	"context"

	"speedtest-tray/internal/speedtest_util"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	WindowWidth  = 320
	WindowHeight = 460
)

type App struct {
	ctx    context.Context
	tester *speedtest_util.SpeedTester
}

func NewApp(tester *speedtest_util.SpeedTester) *App {
	return &App{tester: tester}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) ShowWindow() {
	if a.ctx == nil {
		return
	}

	wailsRuntime.WindowSetSize(a.ctx, WindowWidth, WindowHeight)
	a.positionWindow()
	wailsRuntime.WindowShow(a.ctx)
	wailsRuntime.WindowUnminimise(a.ctx)
	a.ApplyRoundedCorners()
	wailsRuntime.EventsEmit(a.ctx, "window_shown")
}

func (a *App) HideWindow() {
	if a.ctx != nil {
		wailsRuntime.WindowHide(a.ctx)
	}
}

func (a *App) Quit() {
	if a.ctx != nil {
		wailsRuntime.Quit(a.ctx)
	}
}
