package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type AppConfig struct {
	SaveLogs bool `json:"save_logs"`
}

var DefaultConfig = AppConfig{
	SaveLogs: false,
}

// GetAppDir returns the application data directory
// Tries OneDrive/Documents first, falls back to Documents
func GetAppDir() string {
	home, _ := os.UserHomeDir()

	// Try OneDrive first
	dir := filepath.Join(home, "OneDrive", "Documents", AppName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		dir = filepath.Join(home, "Documents", AppName)
	}
	return dir
}

// Load reads config from disk; returns DefaultConfig if not found
func Load() AppConfig {
	dir := GetAppDir()
	configPath := filepath.Join(dir, "config.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return DefaultConfig
	}

	cfg := DefaultConfig
	json.Unmarshal(data, &cfg)
	return cfg
}

// Save writes config to disk
func Save(cfg AppConfig) error {
	dir := GetAppDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("Failed to create config directory: %v\n", err)
		return err
	}

	configPath := filepath.Join(dir, "config.json")
	data, _ := json.MarshalIndent(cfg, "", "  ")
	return os.WriteFile(configPath, data, 0644)
}

// GetLogFilePath returns the path to the log file
func GetLogFilePath() string {
	dir := GetAppDir()
	return filepath.Join(dir, "app.log")
}
