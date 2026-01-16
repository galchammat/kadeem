package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var logger *slog.Logger
var projectRoot string

// QuickfixHandler wraps slog.Handler to produce quickfix-friendly output with file:line
type QuickfixHandler struct {
	out   io.Writer
	level slog.Level
}

func findProjectRoot() string {
	// Start from current working directory
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	// Walk up the directory tree looking for go.mod
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached the root without finding go.mod
			return ""
		}
		dir = parent
	}
}

func NewQuickfixHandler(out io.Writer, level slog.Level) *QuickfixHandler {
	return &QuickfixHandler{out: out, level: level}
}

func (h *QuickfixHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *QuickfixHandler) Handle(_ context.Context, r slog.Record) error {
	// Get caller information
	var file string
	var line int
	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		// Try to get relative path from project root
		if projectRoot != "" {
			if relPath, err := filepath.Rel(projectRoot, f.File); err == nil {
				file = relPath
			} else {
				file = f.File
			}
		} else {
			file = f.File
		}
		line = f.Line
	}

	// Build attributes string as part of the message
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

	// Format: file:line: LEVEL: message attr1=val1, attr2=val2 (timestamp)
	level := r.Level.String()
	msg := r.Message
	timestamp := r.Time.Format("15:04:05.000")

	var output string
	if file != "" {
		if attrs.Len() > 0 {
			output = fmt.Sprintf("%s:%d: %s: %s %s (%s)\n", file, line, level, msg, attrs.String(), timestamp)
		} else {
			output = fmt.Sprintf("%s:%d: %s: %s (%s)\n", file, line, level, msg, timestamp)
		}
	} else {
		if attrs.Len() > 0 {
			output = fmt.Sprintf("%s: %s %s (%s)\n", level, msg, attrs.String(), timestamp)
		} else {
			output = fmt.Sprintf("%s: %s (%s)\n", level, msg, timestamp)
		}
	}

	_, err := h.out.Write([]byte(output))
	return err
}

func (h *QuickfixHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// For simplicity, we don't support persistent attributes in this implementation
	// You could extend this to store attrs and prepend them to each log
	return h
}

func (h *QuickfixHandler) WithGroup(name string) slog.Handler {
	// For simplicity, we don't support groups in this implementation
	return h
}

func init() {
	// Find project root by looking for go.mod
	projectRoot = findProjectRoot()

	// sensible default: quickfix handler, INFO level, stderr
	h := NewQuickfixHandler(os.Stderr, slog.LevelDebug)
	logger = slog.New(h)
	logger = logger.With() // This ensures AddSource works properly
	slog.SetDefault(logger)
}

// Init reconfigures the global logger; call once from main for want custom output/level.
func Init(out io.Writer, level slog.Level) {
	if out == nil {
		out = os.Stderr
	}
	h := NewQuickfixHandler(out, level)
	logger = slog.New(h)
	logger = logger.With() // This ensures AddSource works properly
	slog.SetDefault(logger)
}

// Expose level loggers with caller information
func Info(msg string, args ...any) {
	logWithCaller(slog.LevelInfo, msg, args...)
}

func Debug(msg string, args ...any) {
	logWithCaller(slog.LevelDebug, msg, args...)
}

func Warn(msg string, args ...any) {
	logWithCaller(slog.LevelWarn, msg, args...)
}

func Error(msg string, args ...any) {
	logWithCaller(slog.LevelError, msg, args...)
}

func logWithCaller(level slog.Level, msg string, args ...any) {
	if !logger.Enabled(context.Background(), level) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // skip [Callers, logWithCaller, Info/Debug/Warn/Error]
	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add(args...)
	_ = logger.Handler().Handle(context.Background(), r)
}

// Expose the underlying logger if needed
func Logger() *slog.Logger { return logger }
