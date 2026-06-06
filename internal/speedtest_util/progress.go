package speedtest_util

import "fmt"

// CalculatePhaseProgress calculates progress within a phase given elapsed time and total duration.
// Returns clamped progress (0.0 to 1.0).
func CalculatePhaseProgress(elapsed, duration float64) float64 {
	if duration <= 0 || elapsed <= 0 {
		return 0
	}
	progress := elapsed / duration
	if progress > 1.0 {
		return 1.0
	}
	return progress
}

// MapPhaseProgressToTotal maps phase progress to total test progress.
// Args: phaseStart, phaseEnd = progress range for this phase (0.0-1.0)
//       phaseProgress = progress within phase (0.0-1.0)
// Returns: total test progress
func MapPhaseProgressToTotal(phaseStart, phaseEnd, phaseProgress float64) float64 {
	if phaseProgress < 0 {
		phaseProgress = 0
	}
	if phaseProgress > 1 {
		phaseProgress = 1
	}
	return phaseStart + (phaseProgress * (phaseEnd - phaseStart))
}

// FormatNumber formats a float with specified decimal precision
func FormatNumber(value float64, precision int) string {
	return fmt.Sprintf("%.*f", precision, value)
}
