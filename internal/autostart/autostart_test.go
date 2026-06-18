//go:build !ci

package autostart

import (
	"os"
	"testing"
)

func TestAutostartRoundTrip(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping registry/autostart writes in CI environment")
	}

	mgr, err := New()
	if err != nil {
		t.Fatalf("failed to create Manager: %v", err)
	}

	// Save original status
	originalState := mgr.IsEnabled()

	// Clean up / disable first
	_ = mgr.SetEnabled(false)
	if mgr.IsEnabled() {
		t.Error("expected autostart to be disabled")
	}

	// Enable
	err = mgr.SetEnabled(true)
	if err != nil {
		t.Fatalf("failed to enable autostart: %v", err)
	}
	if !mgr.IsEnabled() {
		t.Error("expected autostart to be enabled")
	}

	// Disable/restore original
	err = mgr.SetEnabled(originalState)
	if err != nil {
		t.Errorf("failed to restore original autostart state: %v", err)
	}
	if mgr.IsEnabled() != originalState {
		t.Errorf("expected restored autostart state to be %v, got %v", originalState, mgr.IsEnabled())
	}
}
