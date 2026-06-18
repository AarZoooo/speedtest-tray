package speedtest_util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHistoryLifecycle(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "speedtest_history_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tempFile := filepath.Join(tempDir, "history.json")
	oldGetHistoryPath := GetHistoryPath
	GetHistoryPath = func() string {
		return tempFile
	}
	defer func() {
		GetHistoryPath = oldGetHistoryPath
	}()

	err = SaveToHistory("Test Server", 15.5, 120.4, 45.2)
	if err != nil {
		t.Fatalf("failed to save to history: %v", err)
	}

	history, err := LoadHistory()
	if err != nil {
		t.Fatalf("failed to load history: %v", err)
	}

	if len(history) != 1 {
		t.Errorf("expected 1 history entry, got %d", len(history))
	}

	entry := history[0]
	if entry.Server != "Test Server" {
		t.Errorf("expected server 'Test Server', got '%s'", entry.Server)
	}
	if entry.Ping != 15.5 {
		t.Errorf("expected ping 15.5, got %f", entry.Ping)
	}
	if entry.Download != 120.4 {
		t.Errorf("expected download 120.4, got %f", entry.Download)
	}
	if entry.Upload != 45.2 {
		t.Errorf("expected upload 45.2, got %f", entry.Upload)
	}

	err = SaveToHistory("Second Server", 12.0, 95.1, 38.0)
	if err != nil {
		t.Fatalf("failed to save second entry: %v", err)
	}

	history, err = LoadHistory()
	if err != nil {
		t.Fatalf("failed to load history: %v", err)
	}

	if len(history) != 2 {
		t.Errorf("expected 2 history entries, got %d", len(history))
	}

	if history[0].Server != "Second Server" {
		t.Errorf("expected newest entry to be 'Second Server', got '%s'", history[0].Server)
	}

	err = ClearHistory()
	if err != nil {
		t.Fatalf("failed to clear history: %v", err)
	}

	history, err = LoadHistory()
	if err != nil {
		t.Fatalf("failed to load cleared history: %v", err)
	}
	if len(history) != 0 {
		t.Errorf("expected 0 history entries, got %d", len(history))
	}
}

func TestHistoryCapping(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "speedtest_history_capping")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tempFile := filepath.Join(tempDir, "history.json")
	oldGetHistoryPath := GetHistoryPath
	GetHistoryPath = func() string {
		return tempFile
	}
	defer func() {
		GetHistoryPath = oldGetHistoryPath
	}()

	for i := 0; i < 60; i++ {
		err := SaveToHistory("Server", 10, 100, 50)
		if err != nil {
			t.Fatalf("failed to save entry %d: %v", i, err)
		}
	}

	history, err := LoadHistory()
	if err != nil {
		t.Fatalf("failed to load history: %v", err)
	}

	if len(history) != 50 {
		t.Errorf("expected history to be capped at 50 entries, got %d", len(history))
	}
}
