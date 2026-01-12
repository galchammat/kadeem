package testlog_test

import (
	"testing"

	"github.com/galchammat/kadeem/internal/testlog"
)

func TestLoggingFormat(t *testing.T) {
	tlog := testlog.New(t)

	tlog.Info("Test started", "testName", "TestLoggingFormat")
	tlog.Debug("Debug information", "key", "value", "count", 42)
	tlog.Warn("Warning message", "component", "testlog")

	// This would fail the test:
	// tlog.Fatalf("This is a fatal error", "error", "something went wrong")

	tlog.Log("Test completed successfully")
}

func TestErrorHandling(t *testing.T) {
	tlog := testlog.New(t)

	tlog.Info("Starting error handling test")

	// Simulate an error condition
	err := someFunction()
	if err != nil {
		// This logs in slog format to both test output and stderr
		tlog.Error("Function failed", "error", err, "function", "someFunction")
		// Test continues...
	}

	tlog.Info("Error handling test completed")
}

func someFunction() error {
	return nil // Simulate success
}
