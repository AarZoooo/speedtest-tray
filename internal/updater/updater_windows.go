package updater

import (
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	"speedtest-tray/internal/config"
)

// Apply downloads the appropriate release asset for the current
// OS/arch, verifies the download size, then performs a platform-
// specific binary swap or installer launch.
func Apply(info UpdateInfo) error {
	if info.DownloadURL == "" {
		return errors.New("empty download URL")
	}

	tempPath := getStagedPath()

	// Download installer
	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Get(info.DownloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(config.ErrUpdateDownload)
	}

	// Verify Content-Length if present
	if resp.ContentLength != -1 && resp.ContentLength != info.AssetSizeBytes {
		return errors.New(config.ErrUpdateBadChecksum)
	}

	out, err := os.Create(tempPath)
	if err != nil {
		return err
	}
	defer out.Close()

	n, err := io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	// Close out before checking/executing
	_ = out.Close()

	// Double check copy size
	if n != info.AssetSizeBytes {
		_ = os.Remove(tempPath)
		return errors.New(config.ErrUpdateBadChecksum)
	}

	// Run Windows installer silently
	cmd := exec.Command(tempPath, "/S")
	if err := cmd.Start(); err != nil {
		return err
	}

	// Exit application
	os.Exit(0)
	return nil
}
