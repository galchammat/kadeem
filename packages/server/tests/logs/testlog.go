package testlog

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// TestHandler is a slog handler that writes to both slog output and testing.T
type TestHandler struct {
	t           *testing.T
	level       slog.Level
	projectRoot string
}

// NewTestHandler creates a handler that logs to both test output and slog
func NewTestHandler(t *testing.T, level slog.Level) *TestHandler {
	projectRoot := findProjectRoot()
	return &TestHandler{
		t:           t,
		level:       level,
		projectRoot: projectRoot,
	}
}

func findProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

func (h *TestHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *TestHandler) Handle(_ context.Context, r slog.Record) error {
	// Get caller information
	var file string
	var line int
	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		if h.projectRoot != "" {
			if relPath, err := filepath.Rel(h.projectRoot, f.File); err == nil {
				file = relPath
			} else {
				file = f.File
			}
		} else {
			file = f.File
		}
		line = f.Line
	}

	// Build attributes string
	var attrs strings.Builder
	r.Attrs(func(a slog.Attr) bool {
		if attrs.Len() > 0 {
			attrs.WriteString(", ")
		}
		attrs.WriteString(a.Key)
		attrs.WriteString("=")
		attrs.WriteString(fmt.Sprint(a.Value))
		return true
	})

	level := r.Level.String()
	msg := r.Message
	timestamp := r.Time.Format("15:04:05.000")

	var output string
	if file != "" {
		if attrs.Len() > 0 {
			output = fmt.Sprintf("%s:%d: %s: %s %s (%s)", file, line, level, msg, attrs.String(), timestamp)
		} else {
			output = fmt.Sprintf("%s:%d: %s: %s (%s)", file, line, level, msg, timestamp)
		}
	} else {
		if attrs.Len() > 0 {
			output = fmt.Sprintf("%s: %s %s (%s)", level, msg, attrs.String(), timestamp)
		} else {
			output = fmt.Sprintf("%s: %s (%s)", level, msg, timestamp)
		}
	}

	// Write to test output
	h.t.Helper()
	h.t.Log(output)

	// Also write to stderr for consistency
	fmt.Fprintln(os.Stderr, output)

	return nil
}

func (h *TestHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *TestHandler) WithGroup(name string) slog.Handler {
	return h
}

// InitTestLogger sets up a logger that writes to both test output and slog
func InitTestLogger(t *testing.T) {
	t.Helper()
	h := NewTestHandler(t, slog.LevelDebug)
	logger := slog.New(h)
	slog.SetDefault(logger)
}

// Helper wraps *testing.T to provide slog-style logging with test control
type Helper struct {
	t *testing.T
}

// New creates a new test helper and initializes the test logger
func New(t *testing.T) *Helper {
	InitTestLogger(t)
	return &Helper{t: t}
}

// logWithCaller logs at the specified level with correct caller information
func (h *Helper) logWithCaller(level slog.Level, msg string, args ...any) {
	h.t.Helper()
	if !slog.Default().Enabled(context.Background(), level) {
		return
	}
	var pcs [1]uintptr
	// Skip 2 frames: [Callers, logWithCaller] to get to the actual test code
	runtime.Callers(3, pcs[:])
	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add(args...)
	_ = slog.Default().Handler().Handle(context.Background(), r)
}

// Fatalf logs an error using slog format and fails the test
func (h *Helper) Fatalf(msg string, args ...any) {
	h.t.Helper()
	h.logWithCaller(slog.LevelError, msg, args...)
	h.t.FailNow()
}

// Fatal logs an error using slog format and fails the test
func (h *Helper) Fatal(msg string, args ...any) {
	h.t.Helper()
	h.logWithCaller(slog.LevelError, msg, args...)
	h.t.FailNow()
}

// Errorf logs an error using slog format but continues the test
func (h *Helper) Errorf(msg string, args ...any) {
	h.t.Helper()
	h.logWithCaller(slog.LevelError, msg, args...)
	h.t.Fail()
}

// Error logs an error using slog format but continues the test
func (h *Helper) Error(msg string, args ...any) {
	h.t.Helper()
	h.logWithCaller(slog.LevelError, msg, args...)
	h.t.Fail()
}

// Logf logs informational messages using slog format
func (h *Helper) Logf(msg string, args ...any) {
	h.t.Helper()
	h.logWithCaller(slog.LevelInfo, msg, args...)
}

// Log logs informational messages using slog format
func (h *Helper) Log(msg string, args ...any) {
	h.t.Helper()
	h.logWithCaller(slog.LevelInfo, msg, args...)
}

// Info is an alias for Log
func (h *Helper) Info(msg string, args ...any) {
	h.t.Helper()
	h.logWithCaller(slog.LevelInfo, msg, args...)
}

// Debug logs debug messages
func (h *Helper) Debug(msg string, args ...any) {
	h.t.Helper()
	h.logWithCaller(slog.LevelDebug, msg, args...)
}

// Warn logs warning messages
func (h *Helper) Warn(msg string, args ...any) {
	h.t.Helper()
	h.logWithCaller(slog.LevelWarn, msg, args...)
}

// T returns the underlying *testing.T
func (h *Helper) T() *testing.T {
	return h.t
}
