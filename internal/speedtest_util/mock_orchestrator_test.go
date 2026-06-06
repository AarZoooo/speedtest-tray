package speedtest_util

import (
	"context"
	"testing"
	"time"
)

var _ TestOrchestrator = (*MockOrchestrator)(nil)

type MockOrchestrator struct {
	GetUserInfoCalls      int
	FindServersCalls      int
	SelectBestServerCalls int
	RunPingCalls          int
	RunDownloadCalls      int
	RunUploadCalls        int

	ServerResult   *ServerInfo
	PingResult     time.Duration
	DownloadResult float64
	UploadResult   float64

	GetUserInfoErr      error
	FindServersErr      error
	SelectBestServerErr error
	RunPingErr          error
	RunDownloadErr      error
	RunUploadErr        error

	DownloadSamples []float64
	UploadSamples   []float64

	OnGetUserInfo      func(context.Context)
	OnFindServers      func(context.Context)
	OnSelectBestServer func(context.Context)
	OnRunPing          func(context.Context)
	OnRunDownload      func(context.Context)
	OnRunUpload        func(context.Context)
}

func (m *MockOrchestrator) GetUserInfo(ctx context.Context) error {
	m.GetUserInfoCalls++
	if m.OnGetUserInfo != nil {
		m.OnGetUserInfo(ctx)
	}
	return m.GetUserInfoErr
}

func (m *MockOrchestrator) FindServers(ctx context.Context) error {
	m.FindServersCalls++
	if m.OnFindServers != nil {
		m.OnFindServers(ctx)
	}
	return m.FindServersErr
}

func (m *MockOrchestrator) SelectBestServer(ctx context.Context) (*ServerInfo, error) {
	m.SelectBestServerCalls++
	if m.OnSelectBestServer != nil {
		m.OnSelectBestServer(ctx)
	}
	if m.SelectBestServerErr != nil {
		return nil, m.SelectBestServerErr
	}
	if m.ServerResult != nil {
		return m.ServerResult, nil
	}
	return &ServerInfo{Name: "Test Server", Country: "Test Country"}, nil
}

func (m *MockOrchestrator) RunPing(ctx context.Context) (time.Duration, error) {
	m.RunPingCalls++
	if m.OnRunPing != nil {
		m.OnRunPing(ctx)
	}
	return m.PingResult, m.RunPingErr
}

func (m *MockOrchestrator) RunDownload(ctx context.Context, callback func(float64)) (float64, error) {
	m.RunDownloadCalls++
	if m.OnRunDownload != nil {
		m.OnRunDownload(ctx)
	}
	for _, sample := range m.DownloadSamples {
		callback(sample)
	}
	return m.DownloadResult, m.RunDownloadErr
}

func (m *MockOrchestrator) RunUpload(ctx context.Context, callback func(float64)) (float64, error) {
	m.RunUploadCalls++
	if m.OnRunUpload != nil {
		m.OnRunUpload(ctx)
	}
	for _, sample := range m.UploadSamples {
		callback(sample)
	}
	return m.UploadResult, m.RunUploadErr
}

func (m *MockOrchestrator) VerifyCalls(t *testing.T, expectedUserInfo, expectedServers, expectedSelect, expectedPing, expectedDL, expectedUL int) {
	t.Helper()
	if m.GetUserInfoCalls != expectedUserInfo {
		t.Errorf("GetUserInfo called %d times, expected %d", m.GetUserInfoCalls, expectedUserInfo)
	}
	if m.FindServersCalls != expectedServers {
		t.Errorf("FindServers called %d times, expected %d", m.FindServersCalls, expectedServers)
	}
	if m.SelectBestServerCalls != expectedSelect {
		t.Errorf("SelectBestServer called %d times, expected %d", m.SelectBestServerCalls, expectedSelect)
	}
	if m.RunPingCalls != expectedPing {
		t.Errorf("RunPing called %d times, expected %d", m.RunPingCalls, expectedPing)
	}
	if m.RunDownloadCalls != expectedDL {
		t.Errorf("RunDownload called %d times, expected %d", m.RunDownloadCalls, expectedDL)
	}
	if m.RunUploadCalls != expectedUL {
		t.Errorf("RunUpload called %d times, expected %d", m.RunUploadCalls, expectedUL)
	}
}
