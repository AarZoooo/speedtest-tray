package gui_wails

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"speedtest-tray/internal/config"
	"speedtest-tray/internal/updater"
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

func TestToggleWindowStress(t *testing.T) {
	app := NewApp()
	app.ctx = context.Background()
	app.isTesting = true

	var wg sync.WaitGroup
	numGoroutines := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			app.ToggleWindow()
		}()
		go func() {
			defer wg.Done()
			app.HideWindow()
		}()
	}

	wg.Wait()
}

func TestGetUpdateInfoReturnsStoredInfo(t *testing.T) {
	app := NewApp()
	want := updater.UpdateInfo{
		LatestVersion:  "1.2.3",
		ReleasePageURL: "https://example.test/releases/v1.2.3",
		AssetSizeBytes: 123,
		HasUpdate:      true,
		DownloadURL:    "https://example.test/download",
	}

	app.mu.Lock()
	app.updateInfo = want
	app.mu.Unlock()

	if got := app.GetUpdateInfo(); got != want {
		t.Fatalf("GetUpdateInfo() = %+v, want %+v", got, want)
	}
}

func TestLaunchAtLoginDefaultsFalseWithoutManager(t *testing.T) {
	app := NewApp()

	if app.GetLaunchAtLogin() {
		t.Fatal("GetLaunchAtLogin() = true, want false without manager")
	}

	app.SetLaunchAtLogin(true)
}

func TestSkipUpdateStoresSkippedVersion(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("APPDATA", dir)
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("HOME", dir)

	app := NewApp()
	app.SkipUpdate("1.2.3")

	got := config.LoadConfigOrDefault()
	if got.SkippedVersion != "1.2.3" {
		t.Fatalf("SkippedVersion = %q, want 1.2.3", got.SkippedVersion)
	}

	configPath := filepath.Join(dir, config.AppName, "config.json")
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("saved config not found at %s: %v", configPath, err)
	}
}
