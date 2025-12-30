package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"
)

var logger *slog.Logger

func init() {
	// sensible default: text handler, INFO level, stderr
	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger = slog.New(h)
	slog.SetDefault(logger)
}

// Init reconfigures the global logger; call once from main for want custom output/level.
func Init(out io.Writer, level slog.Level) {
	if out == nil {
		out = os.Stderr
	}
	h := slog.NewTextHandler(out, &slog.HandlerOptions{
		Level: level,
	})
	logger = slog.New(h)
	slog.SetDefault(logger)
}

// Expose level loggers with correct caller source
func Info(msg string, args ...any) {
	logWithSource(slog.LevelInfo, msg, args...)
}
func Debug(msg string, args ...any) {
	logWithSource(slog.LevelDebug, msg, args...)
}
func Warn(msg string, args ...any) {
	logWithSource(slog.LevelWarn, msg, args...)
}
func Error(msg string, args ...any) {
	logWithSource(slog.LevelError, msg, args...)
}

func logWithSource(level slog.Level, msg string, args ...any) {
	// Manually format the log line
	timestamp := time.Now().Format("2006-01-02T15:04:05.000Z07:00")
	fmt.Fprintf(os.Stderr, "time=%s level=%s ", timestamp, level)

	// Only print call stack for error and warn levels
	if level >= slog.LevelWarn {
		pcs := make([]uintptr, 10)
		n := runtime.Callers(3, pcs) // skip: Callers, logWithSource, Error/Info/etc
		frames := runtime.CallersFrames(pcs[:n])
		var locations []string
		for {
			frame, more := frames.Next()
			if strings.Contains(frame.File, "kadeem") {
				locations = append(locations, fmt.Sprintf("%s:%d", frame.File, frame.Line))
			}
			if !more {
				break
			}
		}
		if len(locations) > 0 {
			fmt.Fprintf(os.Stderr, "%s ", strings.Join(locations, " <- "))
		}
	}

	fmt.Fprintf(os.Stderr, "msg=%q", msg)

	// Print key-value pairs
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			fmt.Fprintf(os.Stderr, " %v=%v", args[i], args[i+1])
		}
	}

	fmt.Fprintln(os.Stderr)
}

// Expose the underlying logger if needed
func Logger() *slog.Logger { return logger }
