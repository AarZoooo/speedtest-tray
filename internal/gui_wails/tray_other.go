//go:build !windows && !darwin

package gui_wails

import (
	"speedtest-tray/internal/config"
)

func StartTray(app *App, iconBytes []byte, macIconBytes []byte, appConfig *config.CustomConfig, toggleLogging func(bool)) {
}

func (a *App) initMacStatusItem() {}

func (a *App) cleanupMacStatusItem() {}
