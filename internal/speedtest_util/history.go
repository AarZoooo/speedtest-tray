package speedtest_util

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"speedtest-tray/internal/config"
)

type HistoryEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Server    string    `json:"server"`
	Ping      float64   `json:"ping"`
	Download  float64   `json:"download"`
	Upload    float64   `json:"upload"`
}

var GetHistoryPath = func() string {
	return filepath.Join(config.GetConfigDir(), "history.json")
}

func SaveToHistory(server string, ping, download, upload float64) error {
	history, err := LoadHistory()
	if err != nil {
		history = []HistoryEntry{}
	}

	newEntry := HistoryEntry{
		Timestamp: time.Now().UTC(),
		Server:    server,
		Ping:      ping,
		Download:  download,
		Upload:    upload,
	}
	history = append([]HistoryEntry{newEntry}, history...)

	if len(history) > config.MaxHistoryEntries {
		history = history[:config.MaxHistoryEntries]
	}

	return saveHistory(history)
}

func LoadHistory() ([]HistoryEntry, error) {
	path := GetHistoryPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var history []HistoryEntry
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, err
	}

	var valid []HistoryEntry
	for _, entry := range history {
		if !entry.Timestamp.IsZero() && (entry.Download > 0 || entry.Upload > 0 || entry.Ping > 0) {
			valid = append(valid, entry)
		}
	}
	return valid, nil
}

func ClearHistory() error {
	return saveHistory([]HistoryEntry{})
}

func saveHistory(history []HistoryEntry) error {
	path := GetHistoryPath()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
