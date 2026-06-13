//go:build !windows && !darwin

package gui_wails

func (a *App) positionWindow() {}

func (a *App) ApplyRoundedCorners() {}
