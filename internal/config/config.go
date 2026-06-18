package config

import (
	"encoding/json"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type CustomConfig struct {
	SaveLogs       bool   `json:"save_logs"`
	SkippedVersion string `json:"skipped_version,omitempty"`
}

var DefaultConfig = CustomConfig{
	SaveLogs:       false,
	SkippedVersion: "",
}

func GetConfigDir() string {
	if IsDev {
		return "."
	}
	configDir, err := os.UserConfigDir()
	if err != nil {
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
		slog.Error(ErrCreateConfigDir, "error", err)
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

var execCommand = exec.Command

func OpenDirectory(path string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = execCommand("explorer", filepath.Clean(path))
	case "darwin":
		cmd = execCommand("open", path)
	default:
		cmd = execCommand("xdg-open", path)
	}
	return cmd.Start()
}
