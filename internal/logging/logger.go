package logging

import (
	"io"
	"log/slog"
	"os"
)

var logger *slog.Logger

func init() {
	// sensible default: text handler, INFO level, stderr
	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger = slog.New(h)
	slog.SetDefault(logger)
}

// Init reconfigures the global logger; call once from main for want custom output/level.
func Init(out io.Writer, level slog.Level) {
	if out == nil {
		out = os.Stderr
	}
	h := slog.NewTextHandler(out, &slog.HandlerOptions{Level: level})
	logger = slog.New(h)
	slog.SetDefault(logger)
}

// Expose level loggers
func Info(msg string, args ...any)  { logger.Info(msg, args...) }
func Debug(msg string, args ...any) { logger.Debug(msg, args...) }
func Warn(msg string, args ...any)  { logger.Warn(msg, args...) }
func Error(msg string, args ...any) { logger.Error(msg, args...) }

// Expose the underlying logger if needed
func Logger() *slog.Logger { return logger }
