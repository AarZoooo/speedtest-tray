package gui_wails

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	tempDir, err := os.MkdirTemp("", "speedtest_gui_wails_test")
	if err != nil {
		panic(err)
	}

	os.Setenv("APPDATA", tempDir)
	os.Setenv("XDG_CONFIG_HOME", tempDir)
	os.Setenv("HOME", tempDir)

	code := m.Run()

	os.RemoveAll(tempDir)
	os.Exit(code)
}
