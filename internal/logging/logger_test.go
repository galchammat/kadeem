package logging

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestInitWithRotation(t *testing.T) {
	// Create a temporary directory for test logs
	tmpDir := t.TempDir()

	// Test with default retention days
	err := InitWithRotation(slog.LevelInfo, tmpDir, 7)
	if err != nil {
		t.Fatalf("Failed to initialize logging with rotation: %v", err)
	}

	// Write some log messages
	Info("Test info message")
	Debug("Test debug message - should not appear with INFO level")
	Warn("Test warning message")
	Error("Test error message")

	// Give it a moment to flush
	time.Sleep(100 * time.Millisecond)

	// Check if log file was created
	logFile := filepath.Join(tmpDir, "kadeem.log")
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Errorf("Log file was not created at %s", logFile)
	}

	// Read the log file and verify it contains our messages
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)
	if len(logContent) == 0 {
		t.Error("Log file is empty")
	}

	// Verify some messages are present (Info and above)
	if !containsString(logContent, "Test info message") {
		t.Error("Log file doesn't contain info message")
	}
	if !containsString(logContent, "Test warning message") {
		t.Error("Log file doesn't contain warning message")
	}
	if !containsString(logContent, "Test error message") {
		t.Error("Log file doesn't contain error message")
	}
}

func TestInitWithRotationFromEnv(t *testing.T) {
	// Create a temporary directory for test logs
	tmpDir := t.TempDir()

	// Test without environment variable (should use default)
	os.Unsetenv("LOG_RETENTION_DAYS")
	err := InitWithRotationFromEnv(slog.LevelDebug, tmpDir)
	if err != nil {
		t.Fatalf("Failed to initialize logging from env: %v", err)
	}

	// Write a test message
	Info("Test message with default retention")
	time.Sleep(100 * time.Millisecond)

	// Check if log file was created
	logFile := filepath.Join(tmpDir, "kadeem.log")
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Errorf("Log file was not created at %s", logFile)
	}

	// Test with environment variable set
	tmpDir2 := t.TempDir()
	os.Setenv("LOG_RETENTION_DAYS", "14")
	defer os.Unsetenv("LOG_RETENTION_DAYS")

	err = InitWithRotationFromEnv(slog.LevelDebug, tmpDir2)
	if err != nil {
		t.Fatalf("Failed to initialize logging from env with custom retention: %v", err)
	}

	Info("Test message with custom retention")
	time.Sleep(100 * time.Millisecond)

	// Check if log file was created
	logFile2 := filepath.Join(tmpDir2, "kadeem.log")
	if _, err := os.Stat(logFile2); os.IsNotExist(err) {
		t.Errorf("Log file was not created at %s", logFile2)
	}
}

func TestInitWithRotationInvalidEnv(t *testing.T) {
	tmpDir := t.TempDir()

	// Test with invalid environment variable (should use default)
	os.Setenv("LOG_RETENTION_DAYS", "invalid")
	defer os.Unsetenv("LOG_RETENTION_DAYS")

	err := InitWithRotationFromEnv(slog.LevelDebug, tmpDir)
	if err != nil {
		t.Fatalf("Failed to initialize logging from env with invalid retention: %v", err)
	}

	Info("Test message with invalid retention value")
	time.Sleep(100 * time.Millisecond)

	// Should still create log file with default retention
	logFile := filepath.Join(tmpDir, "kadeem.log")
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Errorf("Log file was not created at %s", logFile)
	}
}

func containsString(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && s != substr && (len(s) >= len(substr)) && s[:len(substr)] == substr || len(s) > len(substr) && containsHelper(s, substr)
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
