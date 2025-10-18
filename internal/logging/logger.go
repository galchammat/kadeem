package logging

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/natefinch/lumberjack.v2"
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

// InitWithRotation configures the global logger to output to both stdout and rotating log files.
// Log files are rotated daily and old logs are kept for retentionDays.
func InitWithRotation(level slog.Level, logDir string, retentionDays int) error {
	// Create log directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// Configure lumberjack for log rotation
	fileWriter := &lumberjack.Logger{
		Filename:   filepath.Join(logDir, "kadeem.log"),
		MaxSize:    100, // megabytes
		MaxBackups: retentionDays,
		MaxAge:     retentionDays, // days
		Compress:   true,
	}

	// Create a multi-writer that writes to both stdout and the log file
	multiWriter := io.MultiWriter(os.Stdout, fileWriter)

	// Create handler with multi-writer
	h := slog.NewTextHandler(multiWriter, &slog.HandlerOptions{Level: level})
	logger = slog.New(h)
	slog.SetDefault(logger)

	return nil
}

// InitWithRotationFromEnv initializes logging with rotation using environment variables.
// It reads LOG_RETENTION_DAYS from the environment (defaults to 7 if not set or invalid).
func InitWithRotationFromEnv(level slog.Level, logDir string) error {
	retentionDays := 7 // default
	if envDays := os.Getenv("LOG_RETENTION_DAYS"); envDays != "" {
		if days, err := strconv.Atoi(envDays); err == nil && days > 0 {
			retentionDays = days
		}
	}
	return InitWithRotation(level, logDir, retentionDays)
}

// Expose level loggers
func Info(msg string, args ...any)  { logger.Info(msg, args...) }
func Debug(msg string, args ...any) { logger.Debug(msg, args...) }
func Warn(msg string, args ...any)  { logger.Warn(msg, args...) }
func Error(msg string, args ...any) { logger.Error(msg, args...) }

// Expose the underlying logger if needed
func Logger() *slog.Logger { return logger }
