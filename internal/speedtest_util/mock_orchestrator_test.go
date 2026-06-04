package speedtest_util

import (
	"context"
	"testing"
)

// MockOrchestrator is a test implementation of TestOrchestrator
type MockOrchestrator struct {
	// Track call counts
	GetUserInfoCalls      int
	FindServersCalls      int
	SelectBestServerCalls int
	RunPingCalls          int
	RunDownloadCalls      int
	RunUploadCalls        int

	// Control return values
	UserInfoResult *UserInfo
	ServersResult  []Server
	ServerResult   *Server
	PingResult     float64
	DownloadResult float64
	UploadResult   float64
	ErrorResult    error
}

// UserInfo mock types
type UserInfo struct {
	IP   string
	ISP  string
	City string
}

type Server struct {
	ID       string
	Name     string
	Distance float64
}

func (m *MockOrchestrator) GetUserInfo(ctx context.Context) (*UserInfo, error) {
	m.GetUserInfoCalls++
	return m.UserInfoResult, m.ErrorResult
}

func (m *MockOrchestrator) FindServers(ctx context.Context) ([]Server, error) {
	m.FindServersCalls++
	return m.ServersResult, m.ErrorResult
}

func (m *MockOrchestrator) SelectBestServer(ctx context.Context, servers []Server) (*Server, error) {
	m.SelectBestServerCalls++
	return m.ServerResult, m.ErrorResult
}

func (m *MockOrchestrator) RunPing(ctx context.Context, server *Server, callback func(float64)) error {
	m.RunPingCalls++
	if callback != nil && m.PingResult > 0 {
		callback(m.PingResult)
	}
	return m.ErrorResult
}

func (m *MockOrchestrator) RunDownload(ctx context.Context, server *Server, callback func(float64)) error {
	m.RunDownloadCalls++
	if callback != nil && m.DownloadResult > 0 {
		callback(m.DownloadResult)
	}
	return m.ErrorResult
}

func (m *MockOrchestrator) RunUpload(ctx context.Context, server *Server, callback func(float64)) error {
	m.RunUploadCalls++
	if callback != nil && m.UploadResult > 0 {
		callback(m.UploadResult)
	}
	return m.ErrorResult
}

// Test helper to verify mock was called the expected number of times
func (m *MockOrchestrator) VerifyCalls(t *testing.T, expectedUserInfo, expectedServers, expectedSelect, expectedPing, expectedDL, expectedUL int) {
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
