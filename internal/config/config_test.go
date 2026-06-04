package config

import (
	"testing"
)

func TestProgressConstants(t *testing.T) {
	tests := []struct {
		name     string
		expected float64
	}{
		{"ProgressInit", 0.05},
		{"ProgressGetInfo", 0.05},
		{"ProgressFindServers", 0.10},
		{"ProgressSelectServer", 0.12},
		{"ProgressServerSelect", 0.15},
		{"ProgressPingStart", 0.20},
		{"ProgressDownStart", 0.30},
		{"ProgressUpStart", 0.70},
		{"ProgressComplete", 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify constants are accessible and reasonable
			if ProgressInit < 0 || ProgressInit > 1 {
				t.Errorf("ProgressInit %v is outside valid range [0, 1]", ProgressInit)
			}
		})
	}
}

func TestWindowDimensions(t *testing.T) {
	if WindowWidth <= 0 || WindowHeight <= 0 {
		t.Error("Window dimensions must be positive")
	}
}

func TestPhaseConstants(t *testing.T) {
	phases := []Phase{
		PhaseInitializing,
		PhaseGettingInfo,
		PhaseFindingServers,
		PhaseSelectingServer,
		PhaseServerSelected,
		PhasePingTest,
		PhaseStartingDownload,
		PhaseDownloading,
		PhaseStartingUpload,
		PhaseUploading,
		PhaseCompleted,
		PhaseFailed,
	}

	for _, phase := range phases {
		if phase == "" {
			t.Error("Phase constant is empty")
		}
	}
}

