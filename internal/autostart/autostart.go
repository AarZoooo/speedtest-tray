package autostart

import (
	"os"
	"runtime"

	"github.com/emersion/go-autostart"
	"speedtest-tray/internal/config"
)

// Manager wraps go-autostart for the current executable.
type Manager struct {
	app *autostart.App
}

// New creates a Manager pointing at the current executable.
// Returns an error if os.Executable() fails.
func New() (*Manager, error) {
	execPath, err := os.Executable()
	if err != nil {
		return nil, err
	}

	name := config.AppName
	displayName := config.AppName
	if runtime.GOOS == "darwin" {
		name = "dev.aarju.speedtest-tray"
	}

	app := &autostart.App{
		Name:        name,
		DisplayName: displayName,
		Exec:        []string{execPath},
	}

	return &Manager{app: app}, nil
}

// IsEnabled reports whether the app is registered to launch at login.
func (m *Manager) IsEnabled() bool {
	if m.app == nil {
		return false
	}
	return m.app.IsEnabled()
}

// SetEnabled registers or deregisters the app from OS login items.
func (m *Manager) SetEnabled(enabled bool) error {
	if m.app == nil {
		return nil
	}
	if enabled {
		if m.app.IsEnabled() {
			return nil
		}
		return m.app.Enable()
	} else {
		if !m.app.IsEnabled() {
			return nil
		}
		return m.app.Disable()
	}
}
