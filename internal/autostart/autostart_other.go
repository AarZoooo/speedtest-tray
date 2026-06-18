//go:build !windows && !darwin

package autostart

func (m *Manager) IsEnabled() bool {
	return false
}

func (m *Manager) SetEnabled(enabled bool) error {
	return nil
}
