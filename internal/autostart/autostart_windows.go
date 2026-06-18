//go:build windows

package autostart

import "golang.org/x/sys/windows/registry"

func (m *Manager) IsEnabled() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()

	val, _, err := k.GetStringValue(m.name)
	if err != nil {
		return false
	}
	return val == m.execPath
}

func (m *Manager) SetEnabled(enabled bool) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	if enabled {
		return k.SetStringValue(m.name, m.execPath)
	} else {
		err := k.DeleteValue(m.name)
		if err != nil && err != registry.ErrNotExist {
			return err
		}
		return nil
	}
}
