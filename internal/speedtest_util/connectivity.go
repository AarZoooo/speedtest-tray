package speedtest_util

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"syscall"

	"speedtest-tray/internal/config"
)

var checkProbeURL = config.ConnectivityProbeURL

func CheckInternet(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf(config.ErrTestStopped)
	}

	probeCtx, cancel := context.WithTimeout(ctx, config.ConnectivityCheckTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(probeCtx, http.MethodHead, checkProbeURL, nil)
	if err != nil {
		return errors.New(config.ErrNoInternet)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if probeCtx.Err() != nil && errors.Is(ctx.Err(), context.Canceled) {
			return fmt.Errorf(config.ErrTestStopped)
		}
		if isOfflineError(err) {
			return errors.New(config.ErrNoInternet)
		}
		return errors.New(config.ErrNoInternet)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return nil
	}
	return errors.New(config.ErrNoInternet)
}

func isOfflineError(err error) bool {
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}

	var opErr *net.OpError
	if errors.As(err, &opErr) {
		return true
	}

	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return true
	}

	var errno syscall.Errno
	if errors.As(err, &errno) {
		return true
	}

	return false
}
