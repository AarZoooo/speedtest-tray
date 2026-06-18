package updater

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"testing"
)

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		v1, v2 string
		want   int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.0", "1.0.1", -1},
		{"1.0.1", "1.0.0", 1},
		{"v1.2.0", "v1.1.9", 1},
		{"1.2.0", "v1.2.0", 0},
		{"1.10.0", "1.9.0", 1},
	}
	for _, tt := range tests {
		got := compareVersions(tt.v1, tt.v2)
		if got != tt.want {
			t.Errorf("compareVersions(%q, %q) = %d, want %d", tt.v1, tt.v2, got, tt.want)
		}
	}
}

func TestCheck_NewerVersion(t *testing.T) {
	oldOS, oldArch, oldURL := targetOS, targetArch, apiURL
	defer func() {
		targetOS, targetArch, apiURL = oldOS, oldArch, oldURL
	}()

	targetOS = "windows"
	targetArch = "amd64"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		release := githubRelease{
			TagName: "v9.9.9",
			HTMLURL: "https://example.com/release/9.9.9",
			Assets: []githubAsset{
				{
					Name:               "speedtest-tray-windows-amd64.exe",
					Size:               12345,
					BrowserDownloadURL: "https://example.com/download/windows-amd64.exe",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(release)
	}))
	defer server.Close()

	apiURL = server.URL

	info, err := Check("1.0.0", "", "owner", "repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !info.HasUpdate {
		t.Error("expected HasUpdate to be true")
	}
	if info.LatestVersion != "9.9.9" {
		t.Errorf("expected LatestVersion to be 9.9.9, got %s", info.LatestVersion)
	}
	if info.ReleasePageURL != "https://example.com/release/9.9.9" {
		t.Errorf("expected ReleasePageURL to be correct, got %s", info.ReleasePageURL)
	}
	if info.AssetSizeBytes != 12345 {
		t.Errorf("expected AssetSizeBytes to be 12345, got %d", info.AssetSizeBytes)
	}
	if info.DownloadURL != "https://example.com/download/windows-amd64.exe" {
		t.Errorf("expected DownloadURL to be correct, got %s", info.DownloadURL)
	}
}

func TestCheck_SameVersion(t *testing.T) {
	oldOS, oldArch, oldURL := targetOS, targetArch, apiURL
	defer func() {
		targetOS, targetArch, apiURL = oldOS, oldArch, oldURL
	}()

	targetOS = "windows"
	targetArch = "amd64"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		release := githubRelease{
			TagName: "v1.0.0",
			HTMLURL: "https://example.com/release/1.0.0",
			Assets: []githubAsset{
				{
					Name: "speedtest-tray-windows-amd64.exe",
					Size: 12345,
				},
			},
		}
		_ = json.NewEncoder(w).Encode(release)
	}))
	defer server.Close()

	apiURL = server.URL

	info, err := Check("1.0.0", "", "owner", "repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if info.HasUpdate {
		t.Error("expected HasUpdate to be false")
	}
}

func TestCheck_OlderVersion(t *testing.T) {
	oldOS, oldArch, oldURL := targetOS, targetArch, apiURL
	defer func() {
		targetOS, targetArch, apiURL = oldOS, oldArch, oldURL
	}()

	targetOS = "windows"
	targetArch = "amd64"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		release := githubRelease{
			TagName: "v0.9.0",
			HTMLURL: "https://example.com/release/0.9.0",
			Assets: []githubAsset{
				{
					Name: "speedtest-tray-windows-amd64.exe",
					Size: 12345,
				},
			},
		}
		_ = json.NewEncoder(w).Encode(release)
	}))
	defer server.Close()

	apiURL = server.URL

	info, err := Check("1.0.0", "", "owner", "repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if info.HasUpdate {
		t.Error("expected HasUpdate to be false")
	}
}

func TestCheck_SkippedVersion(t *testing.T) {
	oldOS, oldArch, oldURL := targetOS, targetArch, apiURL
	defer func() {
		targetOS, targetArch, apiURL = oldOS, oldArch, oldURL
	}()

	targetOS = "windows"
	targetArch = "amd64"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		release := githubRelease{
			TagName: "v9.9.9",
			HTMLURL: "https://example.com/release/9.9.9",
			Assets: []githubAsset{
				{
					Name: "speedtest-tray-windows-amd64.exe",
					Size: 12345,
				},
			},
		}
		_ = json.NewEncoder(w).Encode(release)
	}))
	defer server.Close()

	apiURL = server.URL

	info, err := Check("1.0.0", "9.9.9", "owner", "repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if info.HasUpdate {
		t.Error("expected HasUpdate to be false")
	}
}

func TestCheck_NetworkError(t *testing.T) {
	oldOS, oldArch, oldURL := targetOS, targetArch, apiURL
	defer func() {
		targetOS, targetArch, apiURL = oldOS, oldArch, oldURL
	}()

	targetOS = "windows"
	targetArch = "amd64"

	// Create and immediately close server to cause a network error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server.Close()

	apiURL = server.URL

	_, err := Check("1.0.0", "", "owner", "repo")
	if err == nil {
		t.Error("expected network error, got nil")
	}
}

func TestCheck_MalformedJSON(t *testing.T) {
	oldOS, oldArch, oldURL := targetOS, targetArch, apiURL
	defer func() {
		targetOS, targetArch, apiURL = oldOS, oldArch, oldURL
	}()

	targetOS = "windows"
	targetArch = "amd64"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{invalid"))
	}))
	defer server.Close()

	apiURL = server.URL

	_, err := Check("1.0.0", "", "owner", "repo")
	if err == nil {
		t.Error("expected json decoding error, got nil")
	}
}

func TestCheck_NoMatchingAsset(t *testing.T) {
	oldOS, oldArch, oldURL := targetOS, targetArch, apiURL
	defer func() {
		targetOS, targetArch, apiURL = oldOS, oldArch, oldURL
	}()

	targetOS = "windows"
	targetArch = "amd64"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		release := githubRelease{
			TagName: "v9.9.9",
			Assets: []githubAsset{
				{
					Name: "speedtest-tray-darwin-amd64.tar.gz",
					Size: 12345,
				},
			},
		}
		_ = json.NewEncoder(w).Encode(release)
	}))
	defer server.Close()

	apiURL = server.URL

	_, err := Check("1.0.0", "", "owner", "repo")
	if err == nil {
		t.Error("expected error for no matching asset, got nil")
	} else if err.Error() != "no matching release asset found for current OS and architecture" {
		t.Errorf("expected no matching asset error message, got: %s", err.Error())
	}
}

func TestAssetSelection_WindowsAmd64(t *testing.T) {
	oldOS, oldArch, oldURL := targetOS, targetArch, apiURL
	defer func() {
		targetOS, targetArch, apiURL = oldOS, oldArch, oldURL
	}()

	targetOS = "windows"
	targetArch = "amd64"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		release := githubRelease{
			TagName: "v2.0.0",
			Assets: []githubAsset{
				{Name: "speedtest-tray-darwin-amd64.tar.gz", Size: 100, BrowserDownloadURL: "url1"},
				{Name: "speedtest-tray-windows-amd64.exe", Size: 200, BrowserDownloadURL: "url2"},
				{Name: "speedtest-tray-windows-arm64.exe", Size: 300, BrowserDownloadURL: "url3"},
			},
		}
		_ = json.NewEncoder(w).Encode(release)
	}))
	defer server.Close()

	apiURL = server.URL
	info, err := Check("1.0.0", "", "owner", "repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.DownloadURL != "url2" {
		t.Errorf("expected url2, got %s", info.DownloadURL)
	}
}

func TestAssetSelection_WindowsArm64(t *testing.T) {
	oldOS, oldArch, oldURL := targetOS, targetArch, apiURL
	defer func() {
		targetOS, targetArch, apiURL = oldOS, oldArch, oldURL
	}()

	targetOS = "windows"
	targetArch = "arm64"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		release := githubRelease{
			TagName: "v2.0.0",
			Assets: []githubAsset{
				{Name: "speedtest-tray-darwin-amd64.tar.gz", Size: 100, BrowserDownloadURL: "url1"},
				{Name: "speedtest-tray-windows-amd64.exe", Size: 200, BrowserDownloadURL: "url2"},
				{Name: "speedtest-tray-windows-arm64.exe", Size: 300, BrowserDownloadURL: "url3"},
			},
		}
		_ = json.NewEncoder(w).Encode(release)
	}))
	defer server.Close()

	apiURL = server.URL
	info, err := Check("1.0.0", "", "owner", "repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.DownloadURL != "url3" {
		t.Errorf("expected url3, got %s", info.DownloadURL)
	}
}

func TestAssetSelection_DarwinAmd64(t *testing.T) {
	oldOS, oldArch, oldURL := targetOS, targetArch, apiURL
	defer func() {
		targetOS, targetArch, apiURL = oldOS, oldArch, oldURL
	}()

	targetOS = "darwin"
	targetArch = "amd64"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		release := githubRelease{
			TagName: "v2.0.0",
			Assets: []githubAsset{
				{Name: "speedtest-tray-darwin-amd64.tar.gz", Size: 100, BrowserDownloadURL: "url1"},
				{Name: "speedtest-tray-windows-amd64.exe", Size: 200, BrowserDownloadURL: "url2"},
				{Name: "speedtest-tray-darwin-arm64.tar.gz", Size: 300, BrowserDownloadURL: "url3"},
			},
		}
		_ = json.NewEncoder(w).Encode(release)
	}))
	defer server.Close()

	apiURL = server.URL
	info, err := Check("1.0.0", "", "owner", "repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.DownloadURL != "url1" {
		t.Errorf("expected url1, got %s", info.DownloadURL)
	}
}

func TestAssetSelection_DarwinArm64(t *testing.T) {
	oldOS, oldArch, oldURL := targetOS, targetArch, apiURL
	defer func() {
		targetOS, targetArch, apiURL = oldOS, oldArch, oldURL
	}()

	targetOS = "darwin"
	targetArch = "arm64"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		release := githubRelease{
			TagName: "v2.0.0",
			Assets: []githubAsset{
				{Name: "speedtest-tray-darwin-amd64.tar.gz", Size: 100, BrowserDownloadURL: "url1"},
				{Name: "speedtest-tray-windows-amd64.exe", Size: 200, BrowserDownloadURL: "url2"},
				{Name: "speedtest-tray-darwin-arm64.tar.gz", Size: 300, BrowserDownloadURL: "url3"},
			},
		}
		_ = json.NewEncoder(w).Encode(release)
	}))
	defer server.Close()

	apiURL = server.URL
	info, err := Check("1.0.0", "", "owner", "repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.DownloadURL != "url3" {
		t.Errorf("expected url3, got %s", info.DownloadURL)
	}
}

func TestCleanupStagedInstaller_FileExists(t *testing.T) {
	oldOS := targetOS
	defer func() { targetOS = oldOS }()
	targetOS = runtime.GOOS

	path := getStagedPath()
	err := os.WriteFile(path, []byte("dummy installer content"), 0644)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	CleanupStagedInstaller()

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected temp file to be cleaned up, but it still exists")
	}
}

func TestCleanupStagedInstaller_NoFile(t *testing.T) {
	oldOS := targetOS
	defer func() { targetOS = oldOS }()
	targetOS = runtime.GOOS

	path := getStagedPath()
	_ = os.Remove(path) // Ensure it doesn't exist

	// Should not panic or error
	CleanupStagedInstaller()
}

func TestApply_ContentLengthMismatch(t *testing.T) {
	if runtime.GOOS != "windows" && runtime.GOOS != "darwin" {
		t.Skip("Apply is only supported on Windows and macOS")
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "100")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(make([]byte, 100))
	}))
	defer server.Close()

	info := UpdateInfo{
		DownloadURL:    server.URL,
		AssetSizeBytes: 50, // Mismatched size
	}

	err := Apply(info, nil)
	if err == nil {
		t.Error("expected error for content length mismatch, got nil")
	} else if err.Error() != "Downloaded installer size mismatch, aborting" {
		t.Errorf("expected size mismatch error, got: %v", err)
	}

	serverChunked := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(make([]byte, 10))
	}))
	defer serverChunked.Close()

	infoChunked := UpdateInfo{
		DownloadURL:    serverChunked.URL,
		AssetSizeBytes: 50,
	}

	err = Apply(infoChunked, nil)
	if err == nil {
		t.Error("expected error for copy size mismatch, got nil")
	} else if err.Error() != "Downloaded installer size mismatch, aborting" {
		t.Errorf("expected size mismatch error, got: %v", err)
	}
}
