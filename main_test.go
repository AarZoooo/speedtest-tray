package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestTruncateLogFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "speedtest-tray-log-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logPath := filepath.Join(tempDir, "test.log")

	truncateLogFile(logPath, 10)

	initialContent := []byte("line 1\nline 2\nline 3\n")
	err = os.WriteFile(logPath, initialContent, 0666)
	if err != nil {
		t.Fatalf("failed to write initial file: %v", err)
	}

	truncateLogFile(logPath, 5)

	readData, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}
	if !bytes.Equal(readData, initialContent) {
		t.Errorf("expected content to be unchanged. got %q, want %q", readData, initialContent)
	}

	truncateLogFile(logPath, 3)
	readData, err = os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}
	if !bytes.Equal(readData, initialContent) {
		t.Errorf("expected content to be unchanged at exact limit. got %q, want %q", readData, initialContent)
	}

	var buffer bytes.Buffer
	for i := 1; i <= 15; i++ {
		buffer.WriteString(fmt.Sprintf("line %d\n", i))
	}
	err = os.WriteFile(logPath, buffer.Bytes(), 0666)
	if err != nil {
		t.Fatalf("failed to write 15 lines: %v", err)
	}

	truncateLogFile(logPath, 10)

	readData, err = os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read truncated log file: %v", err)
	}

	var expectedBuffer bytes.Buffer
	for i := 6; i <= 15; i++ {
		expectedBuffer.WriteString(fmt.Sprintf("line %d\n", i))
	}
	expectedContent := expectedBuffer.Bytes()

	if !bytes.Equal(readData, expectedContent) {
		t.Errorf("incorrect truncation. got:\n%s\nwant:\n%s", readData, expectedContent)
	}
}
