package autostart

import (
	"os"
	"speedtest-tray/internal/config"
)

// Manager wraps autostart configuration for the current executable.
type Manager struct {
	name        string
	displayName string
	execPath    string
}

// New creates a Manager pointing at the current executable.
func New() (*Manager, error) {
	execPath, err := os.Executable()
	if err != nil {
		return nil, err
	}

	name := config.AppName
	displayName := config.AppName

	return &Manager{
		name:        name,
		displayName: displayName,
		execPath:    execPath,
	}, nil
}
