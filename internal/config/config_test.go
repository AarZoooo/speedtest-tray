package config

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestProgressConstants(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		want  float64
	}{
		{"ProgressInit", ProgressInit, 0.00},
		{"ProgressGetInfo", ProgressGetInfo, 0.05},
		{"ProgressFindServers", ProgressFindServers, 0.10},
		{"ProgressSelectServer", ProgressSelectServer, 0.12},
		{"ProgressServerSelect", ProgressServerSelect, 0.15},
		{"ProgressPingStart", ProgressPingStart, 0.20},
		{"ProgressPingEnd", ProgressPingEnd, 0.30},
		{"ProgressDownStart", ProgressDownStart, 0.30},
		{"ProgressDownEnd", ProgressDownEnd, 0.70},
		{"ProgressUpStart", ProgressUpStart, 0.70},
		{"ProgressUpEnd", ProgressUpEnd, 0.95},
		{"ProgressComplete", ProgressComplete, 1.0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.value != test.want {
				t.Fatalf("%s = %v, want %v", test.name, test.value, test.want)
			}
			if test.value < 0 || test.value > 1 {
				t.Fatalf("%s = %v, want value in [0, 1]", test.name, test.value)
			}
		})
	}
}

func TestProgressConstantsAreMonotonic(t *testing.T) {
	values := []float64{
		ProgressInit,
		ProgressGetInfo,
		ProgressFindServers,
		ProgressSelectServer,
		ProgressServerSelect,
		ProgressPingStart,
		ProgressPingEnd,
		ProgressDownStart,
		ProgressDownEnd,
		ProgressUpStart,
		ProgressUpEnd,
		ProgressComplete,
	}

	for i := 1; i < len(values); i++ {
		if values[i] < values[i-1] {
			t.Fatalf("progress value at index %d decreased from %v to %v", i, values[i-1], values[i])
		}
	}
}

func TestTimingConstants(t *testing.T) {
	tests := []struct {
		name  string
		value time.Duration
	}{
		{"EstimatedDurationDownload", EstimatedDurationDownload},
		{"EstimatedDurationUpload", EstimatedDurationUpload},
		{"ResultTimeout", ResultTimeout},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.value <= 0 {
				t.Fatalf("%s must be positive", test.name)
			}
		})
	}

	if UIHideDelayMs <= 0 {
		t.Fatal("UIHideDelayMs must be positive")
	}
}

func TestWindowAndGaugeConstants(t *testing.T) {
	if WindowWidth <= 0 || WindowHeight <= 0 {
		t.Fatal("window dimensions must be positive")
	}
	if WindowCornerRadius < 0 {
		t.Fatal("window corner radius cannot be negative")
	}
	if GaugeMaxDownload <= 0 || GaugeMaxUpload <= 0 {
		t.Fatal("gauge scales must be positive")
	}
}

func TestPhaseConstantsAndStrings(t *testing.T) {
	tests := []struct {
		phase Phase
		want  string
	}{
		{PhaseInitializing, "Initializing..."},
		{PhaseGettingInfo, "Fetching user info..."},
		{PhaseFindingServers, "Locating servers..."},
		{PhaseSelectingServer, "Selecting best server..."},
		{PhaseServerSelected, "Server selected"},
		{PhasePingTest, "Measuring latency..."},
		{PhaseStartingDownload, "Starting download test..."},
		{PhaseDownloading, "Running download test..."},
		{PhaseStartingUpload, "Starting upload test..."},
		{PhaseUploading, "Running upload test..."},
		{PhaseCompleted, "Test Completed"},
		{PhaseFailed, "Test Failed"},
		{PhaseStopped, "Test Stopped"},
		{Phase("CUSTOM"), "CUSTOM"},
	}

	for _, test := range tests {
		t.Run(string(test.phase), func(t *testing.T) {
			if test.phase == "" {
				t.Fatal("phase constant is empty")
			}
			if got := test.phase.String(); got != test.want {
				t.Fatalf("Phase.String() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestLoadConfigOrDefaultMissingFile(t *testing.T) {
	setConfigHome(t)

	if got := LoadConfigOrDefault(); got != DefaultConfig {
		t.Fatalf("LoadConfigOrDefault() = %+v, want %+v", got, DefaultConfig)
	}
}

func TestSaveConfigAndLoadConfigOrDefault(t *testing.T) {
	dir := setConfigHome(t)
	want := CustomConfig{SaveLogs: true}

	if err := SaveConfig(want); err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	configPath := filepath.Join(dir, AppName, "config.json")
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("saved config not found at %s: %v", configPath, err)
	}

	got := LoadConfigOrDefault()
	if got != want {
		t.Fatalf("LoadConfigOrDefault() = %+v, want %+v", got, want)
	}
}

func TestLoadConfigOrDefaultMalformedJSON(t *testing.T) {
	dir := setConfigHome(t)
	configDir := filepath.Join(dir, AppName)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "config.json"), []byte("{"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if got := LoadConfigOrDefault(); got != DefaultConfig {
		t.Fatalf("LoadConfigOrDefault() = %+v, want %+v", got, DefaultConfig)
	}
}

func TestGetLogFilePathUsesConfigDir(t *testing.T) {
	dir := setConfigHome(t)
	want := filepath.Join(dir, AppName, "app.log")

	if got := GetLogFilePath(); got != want {
		t.Fatalf("GetLogFilePath() = %q, want %q", got, want)
	}
}

func TestSaveConfigWritesJSON(t *testing.T) {
	dir := setConfigHome(t)
	want := CustomConfig{SaveLogs: true}

	if err := SaveConfig(want); err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, AppName, "config.json"))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	var got CustomConfig
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if got != want {
		t.Fatalf("saved JSON = %+v, want %+v", got, want)
	}
}

func setConfigHome(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("APPDATA", dir)
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("HOME", dir)

	configDir, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("UserConfigDir error = %v", err)
	}
	return configDir
}

func TestOpenDirectory(t *testing.T) {
	var calledName string
	var calledArgs []string

	oldExecCommand := execCommand
	defer func() { execCommand = oldExecCommand }()

	execCommand = func(name string, arg ...string) *exec.Cmd {
		calledName = name
		calledArgs = arg
		return oldExecCommand("go", "version")
	}

	err := OpenDirectory("/test/path")
	if err != nil {
		t.Fatalf("OpenDirectory() error = %v", err)
	}

	var expectedCmd string
	switch runtime.GOOS {
	case "windows":
		expectedCmd = "explorer"
	case "darwin":
		expectedCmd = "open"
	default:
		expectedCmd = "xdg-open"
	}

	if calledName != expectedCmd {
		t.Errorf("expected command %q, got %q", expectedCmd, calledName)
	}

	if len(calledArgs) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(calledArgs))
	}

	expectedPath := "/test/path"
	if runtime.GOOS == "windows" {
		expectedPath = filepath.Clean(expectedPath)
	}
	if calledArgs[0] != expectedPath {
		t.Errorf("expected arg %q, got %q", expectedPath, calledArgs[0])
	}
}
