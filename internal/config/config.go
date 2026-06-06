package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type CustomConfig struct {
	SaveLogs bool `json:"save_logs"`
}

var DefaultConfig = CustomConfig{
	SaveLogs: false,
}

func GetConfigDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		// Fallback to home directory if config dir is unavailable
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "."+filepath.Base(AppName))
	}
	return filepath.Join(configDir, AppName)
}

func LoadConfigOrDefault() CustomConfig {
	dir := GetConfigDir()
	configPath := filepath.Join(dir, "config.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return DefaultConfig
	}

	cfg := DefaultConfig
	json.Unmarshal(data, &cfg)
	return cfg
}

func SaveConfig(cfg CustomConfig) error {
	dir := GetConfigDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("Failed to create config directory: %v\n", err)
		return err
	}

	configPath := filepath.Join(dir, "config.json")
	data, _ := json.MarshalIndent(cfg, "", "  ")
	return os.WriteFile(configPath, data, 0644)
}

func GetLogFilePath() string {
	dir := GetConfigDir()
	return filepath.Join(dir, "app.log")
}
