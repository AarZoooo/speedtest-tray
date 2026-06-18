package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"speedtest-tray/internal/config"
	"speedtest-tray/internal/speedtest_util"
)

type mockOrchestrator struct {
	serverErr error
	pingErr   error
	dlErr     error
	ulErr     error
}

func (m *mockOrchestrator) GetUserInfo(ctx context.Context) error {
	return nil
}

func (m *mockOrchestrator) FindServers(ctx context.Context) error {
	return nil
}

func (m *mockOrchestrator) SelectBestServer(ctx context.Context) (*speedtest_util.ServerInfo, error) {
	if m.serverErr != nil {
		return nil, m.serverErr
	}
	return &speedtest_util.ServerInfo{Name: "Mock Server", Country: "Mock Country"}, nil
}

func (m *mockOrchestrator) RunPing(ctx context.Context) (time.Duration, error) {
	if m.pingErr != nil {
		return 0, m.pingErr
	}
	return 10 * time.Millisecond, nil
}

func (m *mockOrchestrator) RunDownload(ctx context.Context, callback func(float64)) (float64, error) {
	if m.dlErr != nil {
		return 0, m.dlErr
	}
	callback(100.5)
	return 100.5, nil
}

func (m *mockOrchestrator) RunUpload(ctx context.Context, callback func(float64)) (float64, error) {
	if m.ulErr != nil {
		return 0, m.ulErr
	}
	callback(50.2)
	return 50.2, nil
}

func TestRunPrettySuccess(t *testing.T) {
	mock := &mockOrchestrator{}
	buf := &bytes.Buffer{}
	ctx := context.Background()

	err := Run(ctx, false, "", mock, buf)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, config.CLIHeader) {
		t.Errorf("expected output to contain CLI header")
	}
	if !strings.Contains(out, "Mock Server") {
		t.Errorf("expected output to contain server name")
	}
	if !strings.Contains(out, "10 ms") && !strings.Contains(out, "10.00 ms") {
		t.Errorf("expected output to contain ping result")
	}
	if !strings.Contains(out, "100.50 Mbps") {
		t.Errorf("expected output to contain download speed")
	}
	if !strings.Contains(out, "50.20 Mbps") {
		t.Errorf("expected output to contain upload speed")
	}
}

func TestRunPrettyFailure(t *testing.T) {
	mock := &mockOrchestrator{
		serverErr: errors.New("lookup failed"),
	}
	buf := &bytes.Buffer{}
	ctx := context.Background()

	err := Run(ctx, false, "", mock, buf)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	out := buf.String()
	if !strings.Contains(out, "lookup failed") {
		t.Errorf("expected output to contain error message, got %q", out)
	}
}

func TestRunJSONSuccess(t *testing.T) {
	mock := &mockOrchestrator{}
	buf := &bytes.Buffer{}
	ctx := context.Background()

	err := Run(ctx, true, "", mock, buf)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var output jsonOutput
	if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if output.Status != config.JSONStatusSuccess {
		t.Errorf("expected status %q, got %q", config.JSONStatusSuccess, output.Status)
	}
	if output.PingMS != 10 {
		t.Errorf("expected ping 10, got %f", output.PingMS)
	}
	if output.DownloadMbps != 100.5 {
		t.Errorf("expected download 100.5, got %f", output.DownloadMbps)
	}
	if output.UploadMbps != 50.2 {
		t.Errorf("expected upload 50.2, got %f", output.UploadMbps)
	}
	if output.Server != "Mock Server (Mock Country)" {
		t.Errorf("expected server name, got %q", output.Server)
	}
}

func TestRunJSONFailure(t *testing.T) {
	mock := &mockOrchestrator{
		serverErr: errors.New("no connection"),
	}
	buf := &bytes.Buffer{}
	ctx := context.Background()

	err := Run(ctx, true, "", mock, buf)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var output jsonOutput
	if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if output.Status != config.JSONStatusFailed {
		t.Errorf("expected status %q, got %q", config.JSONStatusFailed, output.Status)
	}
	if output.Error != "no connection" {
		t.Errorf("expected error %q, got %q", "no connection", output.Error)
	}
}

func TestPrintHistory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "speedtest_cli_history_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tempFile := filepath.Join(tempDir, "history.json")
	oldGetHistoryPath := speedtest_util.GetHistoryPath
	speedtest_util.GetHistoryPath = func() string {
		return tempFile
	}
	defer func() {
		speedtest_util.GetHistoryPath = oldGetHistoryPath
	}()

	buf := &bytes.Buffer{}
	err = PrintHistory(buf, false)
	if err != nil {
		t.Fatalf("PrintHistory empty failed: %v", err)
	}
	if !strings.Contains(buf.String(), "No speedtest history found.") {
		t.Errorf("expected empty history message, got %q", buf.String())
	}

	err = speedtest_util.SaveToHistory("Server A", 10.0, 100.5, 50.2)
	if err != nil {
		t.Fatalf("failed to save run: %v", err)
	}

	buf.Reset()
	err = PrintHistory(buf, false)
	if err != nil {
		t.Fatalf("PrintHistory tabular failed: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Speedtest History:") {
		t.Errorf("expected header 'Speedtest History:', got %q", out)
	}
	if !strings.Contains(out, "Server A") {
		t.Errorf("expected output to contain server name, got %q", out)
	}
	if !strings.Contains(out, "100.5 Mbps") {
		t.Errorf("expected output to contain download speed, got %q", out)
	}

	buf.Reset()
	err = PrintHistory(buf, true)
	if err != nil {
		t.Fatalf("PrintHistory json failed: %v", err)
	}
	var history []speedtest_util.HistoryEntry
	if err := json.Unmarshal(buf.Bytes(), &history); err != nil {
		t.Fatalf("failed to parse json history: %v", err)
	}
	if len(history) != 1 || history[0].Server != "Server A" {
		t.Errorf("unexpected json history: %v", history)
	}
}
