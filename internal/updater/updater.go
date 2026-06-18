package updater

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const updaterUserAgent = "speedtest-tray-updater"

var (
	apiURL     = "https://api.github.com"
	targetOS   = runtime.GOOS
	targetArch = runtime.GOARCH
)

// UpdateInfo holds the result of a version check.
type UpdateInfo struct {
	LatestVersion  string
	ReleasePageURL string
	AssetSizeBytes int64
	HasUpdate      bool
	DownloadURL    string
}

type githubAsset struct {
	Name               string `json:"name"`
	Size               int64  `json:"size"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type githubRelease struct {
	TagName string        `json:"tag_name"`
	HTMLURL string        `json:"html_url"`
	Assets  []githubAsset `json:"assets"`
}

// compareVersions returns:
// -1 if v1 < v2
//  0 if v1 == v2
//  1 if v1 > v2
func compareVersions(v1, v2 string) int {
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")
	for i := 0; i < len(parts1) || i < len(parts2); i++ {
		var n1, n2 int
		if i < len(parts1) {
			n1, _ = strconv.Atoi(parts1[i])
		}
		if i < len(parts2) {
			n2, _ = strconv.Atoi(parts2[i])
		}
		if n1 < n2 {
			return -1
		} else if n1 > n2 {
			return 1
		}
	}
	return 0
}

// Check calls the GitHub Releases API and returns update info.
// Returns an UpdateInfo with HasUpdate=false if the current version
// is up to date or the skipped version matches the latest.
func Check(currentVersion, skippedVersion, owner, repo string) (UpdateInfo, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/releases/latest", apiURL, owner, repo)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return UpdateInfo{}, err
	}
	req.Header.Set("User-Agent", updaterUserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return UpdateInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return UpdateInfo{}, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return UpdateInfo{}, err
	}

	latestVersion := strings.TrimPrefix(release.TagName, "v")
	if latestVersion == skippedVersion {
		return UpdateInfo{
			LatestVersion:  latestVersion,
			ReleasePageURL: release.HTMLURL,
			HasUpdate:      false,
		}, nil
	}

	if compareVersions(latestVersion, currentVersion) <= 0 {
		return UpdateInfo{
			LatestVersion:  latestVersion,
			ReleasePageURL: release.HTMLURL,
			HasUpdate:      false,
		}, nil
	}

	// Find the appropriate asset
	var targetAsset *githubAsset
	for _, asset := range release.Assets {
		lowerName := strings.ToLower(asset.Name)
		var matchesOS bool
		var matchesExt bool

		if targetOS == "windows" {
			matchesOS = strings.Contains(lowerName, "windows") && strings.Contains(lowerName, targetArch)
			matchesExt = strings.HasSuffix(lowerName, ".exe") || strings.HasSuffix(lowerName, ".zip")
		} else if targetOS == "darwin" {
			matchesOS = strings.Contains(lowerName, "darwin") && strings.Contains(lowerName, targetArch)
			matchesExt = strings.HasSuffix(lowerName, ".tar.gz")
		}

		if matchesOS && matchesExt {
			targetAsset = &asset
			break
		}
	}

	if targetAsset == nil {
		return UpdateInfo{}, errors.New("no matching release asset found for current OS and architecture")
	}

	return UpdateInfo{
		LatestVersion:  latestVersion,
		ReleasePageURL: release.HTMLURL,
		AssetSizeBytes: targetAsset.Size,
		HasUpdate:      true,
		DownloadURL:    targetAsset.BrowserDownloadURL,
	}, nil
}

func getStagedPath() string {
	if targetOS == "windows" {
		return filepath.Join(os.TempDir(), "speedtest-tray-update.exe")
	}
	return filepath.Join(os.TempDir(), "speedtest-tray-update")
}

// CleanupStagedInstaller removes any leftover temp installer file
// from a previous update attempt. Safe to call even if no file exists.
func CleanupStagedInstaller() {
	path := getStagedPath()
	if _, err := os.Stat(path); err == nil {
		_ = os.Remove(path)
	}
}

type progressWriter struct {
	w          io.Writer
	total      int64
	written    int64
	onProgress func(percent int)
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n, err := pw.w.Write(p)
	if err != nil {
		return n, err
	}
	pw.written += int64(n)
	if pw.total > 0 && pw.onProgress != nil {
		percent := int(pw.written * 100 / pw.total)
		pw.onProgress(percent)
	}
	return n, nil
}
