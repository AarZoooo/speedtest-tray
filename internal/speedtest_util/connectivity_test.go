package speedtest_util

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"speedtest-tray/internal/config"
)

func TestCheckInternetSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	restore := overrideProbeURL(server.URL)
	defer restore()

	if err := CheckInternet(context.Background()); err != nil {
		t.Fatalf("CheckInternet() error = %v, want nil", err)
	}
}

func TestCheckInternetOffline(t *testing.T) {
	restore := overrideProbeURL("http://127.0.0.1:1")
	defer restore()

	err := CheckInternet(context.Background())
	if err == nil || err.Error() != config.ErrNoInternet {
		t.Fatalf("CheckInternet() error = %v, want %q", err, config.ErrNoInternet)
	}
}

func TestCheckInternetCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := CheckInternet(ctx)
	if err == nil || err.Error() != config.ErrTestStopped {
		t.Fatalf("CheckInternet() error = %v, want %q", err, config.ErrTestStopped)
	}
}

func TestIsOfflineError(t *testing.T) {
	if !isOfflineError(&net.DNSError{IsNotFound: true}) {
		t.Fatal("expected DNS error to be offline")
	}
	if !isOfflineError(&net.OpError{Op: "dial", Net: "tcp", Err: errors.New("connection refused")}) {
		t.Fatal("expected op error to be offline")
	}
	if isOfflineError(nil) {
		t.Fatal("nil should not be offline")
	}
}

func overrideProbeURL(url string) func() {
	original := checkProbeURL
	checkProbeURL = url
	return func() {
		checkProbeURL = original
	}
}
