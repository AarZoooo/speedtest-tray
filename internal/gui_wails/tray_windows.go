//go:build windows

package gui_wails

import (
	"github.com/energye/systray"
	"speedtest-tray/internal/config"
)

func StartTray(app *App, iconBytes []byte, macIconBytes []byte, appConfig *config.CustomConfig, toggleLogging func(bool)) {
	go systray.Run(func() {
		systray.SetTitle(config.AppName)
		systray.SetTooltip(config.AppName)
		systray.SetTemplateIcon(macIconBytes, iconBytes)
		systray.SetOnClick(func(menu systray.IMenu) {
			go app.ToggleWindow()
		})

		show := systray.AddMenuItem("Show", "Show the speedtest window")
		show.Click(app.ShowWindow)

		systray.AddSeparator()

		saveLogs := systray.AddMenuItemCheckbox("Enable Session Logging", "Save test logs to app data", appConfig.SaveLogs)
		saveLogs.Click(func() {
			if saveLogs.Checked() {
				saveLogs.Uncheck()
				toggleLogging(false)
			} else {
				saveLogs.Check()
				toggleLogging(true)
			}
		})

		systray.AddSeparator()

		quit := systray.AddMenuItem("Quit", "Quit the application")
		quit.Click(func() {
			app.Quit()
			systray.Quit()
		})
	}, func() {})
}

func (a *App) initMacStatusItem() {}

func (a *App) cleanupMacStatusItem() {}
