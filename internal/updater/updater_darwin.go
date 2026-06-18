package updater

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"speedtest-tray/internal/config"
)

// Apply downloads the appropriate release asset for the current
// OS/arch, verifies the download size, then performs a platform-
// specific binary swap or installer launch.
func Apply(info UpdateInfo, onProgress func(percent int)) error {
	if info.DownloadURL == "" {
		return errors.New("empty download URL")
	}

	tempTarPath := getStagedPath() + ".tar.gz"
	defer func() {
		_ = os.Remove(tempTarPath)
	}()

	// Download tar.gz
	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Get(info.DownloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(config.ErrUpdateDownload)
	}

	// Verify Content-Length
	if resp.ContentLength != -1 && resp.ContentLength != info.AssetSizeBytes {
		return errors.New(config.ErrUpdateBadChecksum)
	}

	out, err := os.Create(tempTarPath)
	if err != nil {
		return err
	}

	pw := &progressWriter{
		w:          out,
		total:      info.AssetSizeBytes,
		onProgress: onProgress,
	}

	n, err := io.Copy(pw, resp.Body)
	_ = out.Close()
	if err != nil {
		return err
	}

	if n != info.AssetSizeBytes {
		return errors.New(config.ErrUpdateBadChecksum)
	}

	// Extract binary from tar.gz
	archiveFile, err := os.Open(tempTarPath)
	if err != nil {
		return err
	}
	defer archiveFile.Close()

	gzipReader, err := gzip.NewReader(archiveFile)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	extractedBinaryPath := getStagedPath()
	found := false

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Look for the binary
		if header.Typeflag == tar.TypeReg && filepath.Base(header.Name) == config.AppName {
			binOut, err := os.OpenFile(extractedBinaryPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
			if err != nil {
				return err
			}
			_, err = io.Copy(binOut, tarReader)
			_ = binOut.Close()
			if err != nil {
				return err
			}
			found = true
			break
		}
	}

	if !found {
		return errors.New("binary not found in release archive")
	}
	defer func() {
		_ = os.Remove(extractedBinaryPath)
	}()

	// Remove macOS quarantine attribute
	_ = exec.Command("xattr", "-d", "com.apple.quarantine", extractedBinaryPath).Run()

	// Get current executable path
	currentExec, err := os.Executable()
	if err != nil {
		return err
	}

	// Try simple rename/replace
	oldPath := currentExec + ".old"
	_ = os.Remove(oldPath)

	err = os.Rename(currentExec, oldPath)
	if err == nil {
		err = os.Rename(extractedBinaryPath, currentExec)
		if err != nil {
			// Rollback
			_ = os.Rename(oldPath, currentExec)
			return err
		}
		_ = os.Remove(oldPath)
	} else {
		// On permission/EPERM error, use osascript to escalate privileges
		script := "do shell script \"cp -f '" + extractedBinaryPath + "' '" + currentExec + "'\" with administrator privileges"
		cmd := exec.Command("osascript", "-e", script)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	// Relaunch the new binary
	err = syscall.Exec(currentExec, os.Args, os.Environ())
	if err != nil {
		// Fallback: start process and exit
		cmd := exec.Command(currentExec, os.Args[1:]...)
		if err := cmd.Start(); err != nil {
			return err
		}
		os.Exit(0)
	}

	return nil
}
