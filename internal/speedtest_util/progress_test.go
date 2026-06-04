package speedtest_util

import (
	"testing"
)

func TestCalculatePhaseProgress(t *testing.T) {
	tests := []struct {
		name        string
		elapsed     float64
		duration    float64
		expected    float64
		description string
	}{
		{
			name:        "complete_phase",
			elapsed:     10,
			duration:    10,
			expected:    1.0,
			description: "Phase complete",
		},
		{
			name:        "half_phase",
			elapsed:     5,
			duration:    10,
			expected:    0.5,
			description: "Phase halfway",
		},
		{
			name:        "zero_progress",
			elapsed:     0,
			duration:    10,
			expected:    0.0,
			description: "No elapsed time",
		},
		{
			name:        "exceed_duration",
			elapsed:     15,
			duration:    10,
			expected:    1.0,
			description: "Elapsed exceeds duration, clamped to 1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculatePhaseProgress(tt.elapsed, tt.duration)
			if result != tt.expected {
				t.Errorf("CalculatePhaseProgress(%v, %v) = %v, want %v",
					tt.elapsed, tt.duration, result, tt.expected)
			}
		})
	}
}

func TestMapPhaseProgressToTotal(t *testing.T) {
	tests := []struct {
		name          string
		phaseStart    float64
		phaseEnd      float64
		phaseProgress float64
		expected      float64
		description   string
	}{
		{
			name:          "download_phase_zero",
			phaseStart:    0.5,
			phaseEnd:      0.9,
			phaseProgress: 0.0,
			expected:      0.5,
			description:   "Download phase at start",
		},
		{
			name:          "download_phase_full",
			phaseStart:    0.5,
			phaseEnd:      0.9,
			phaseProgress: 1.0,
			expected:      0.9,
			description:   "Download phase complete",
		},
		{
			name:          "download_phase_half",
			phaseStart:    0.5,
			phaseEnd:      0.9,
			phaseProgress: 0.5,
			expected:      0.7,
			description:   "Download phase halfway (0.5 + (0.9-0.5)*0.5 = 0.7)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapPhaseProgressToTotal(tt.phaseStart, tt.phaseEnd, tt.phaseProgress)
			if result != tt.expected {
				t.Errorf("MapPhaseProgressToTotal(%v, %v, %v) = %v, want %v",
					tt.phaseStart, tt.phaseEnd, tt.phaseProgress, result, tt.expected)
			}
		})
	}
}

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		name      string
		input     float64
		precision int
		expected  string
	}{
		{"zero", 0, 2, "0.00"},
		{"integer", 42, 2, "42.00"},
		{"one_decimal", 42.1, 2, "42.10"},
		{"two_decimals", 42.12, 2, "42.12"},
		{"three_decimals_rounded", 42.126, 2, "42.13"},
		{"small_value", 0.001, 2, "0.00"},
		{"large_value", 1234.567, 2, "1234.57"},
		{"single_precision", 42.5, 1, "42.5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatNumber(tt.input, tt.precision)
			if result != tt.expected {
				t.Errorf("FormatNumber(%v, %v) = %q, want %q", tt.input, tt.precision, result, tt.expected)
			}
		})
	}
}
