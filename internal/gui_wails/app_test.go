package gui_wails

import (
	"context"
	"testing"
	"time"

	"speedtest-tray/internal/config"
)

func TestToggleWindowThreshold(t *testing.T) {
	app := NewApp()
	app.ctx = context.Background()
	app.isTesting = true

	if app.windowVisible {
		t.Fatal("Expected initial windowVisible to be false")
	}

	app.ToggleWindow()
	if !app.windowVisible {
		t.Fatal("Expected windowVisible to be true after first toggle")
	}

	app.ToggleWindow()
	if app.windowVisible {
		t.Fatal("Expected windowVisible to be false after second toggle")
	}
	if app.lastHiddenTime.IsZero() {
		t.Fatal("Expected lastHiddenTime to be set")
	}

	app.ToggleWindow()
	if app.windowVisible {
		t.Fatal("Expected ToggleWindow to be ignored immediately after hiding")
	}

	time.Sleep(config.ToggleThreshold + 50*time.Millisecond)
	app.ToggleWindow()
	if !app.windowVisible {
		t.Fatal("Expected windowVisible to be true after waiting and toggling")
	}
}
