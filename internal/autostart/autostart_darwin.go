//go:build darwin

package autostart

import "github.com/emersion/go-autostart"

func (m *Manager) getApp() *autostart.App {
	return &autostart.App{
		Name:        "dev.aarju.speedtest-tray",
		DisplayName: m.displayName,
		Exec:        []string{m.execPath},
	}
}

func (m *Manager) IsEnabled() bool {
	return m.getApp().IsEnabled()
}

func (m *Manager) SetEnabled(enabled bool) error {
	app := m.getApp()
	if enabled {
		if app.IsEnabled() {
			return nil
		}
		return app.Enable()
	} else {
		if !app.IsEnabled() {
			return nil
		}
		return app.Disable()
	}
}
