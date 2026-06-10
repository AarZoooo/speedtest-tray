package speedtest_util

import "fmt"

func CalculatePhaseProgress(elapsed, duration float64) float64 {
	if duration <= 0 || elapsed <= 0 {
		return 0
	}

	if elapsed >= duration {
		return 1.0 - 0.01*(duration/elapsed)
	}

	return elapsed / duration
}

func MapPhaseProgressToTotal(phaseStart, phaseEnd, phaseProgress float64) float64 {
	if phaseProgress < 0 {
		phaseProgress = 0
	}
	if phaseProgress > 1 {
		phaseProgress = 1
	}
	return phaseStart + (phaseProgress * (phaseEnd - phaseStart))
}

func FormatNumber(value float64, precision int) string {
	return fmt.Sprintf("%.*f", precision, value)
}
