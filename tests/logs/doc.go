package testlog

/*
Package testlog provides a slog-compatible test logging helper that formats
test output in the same quickfix-friendly format as the main application logger.

Usage:

Instead of using t.Fatalf with printf-style formatting:
    t.Fatalf("Failed to open database: %v", err)

Use testlog with slog-style key-value pairs:
    tlog := testlog.New(t)
    tlog.Fatalf("Failed to open database", "error", err)

This ensures all log output, including test failures, appears in the same
format with file:line information and timestamps.

Example:

    func TestDatabase(t *testing.T) {
        tlog := testlog.New(t)

        db, err := database.Open()
        if err != nil {
            tlog.Fatalf("Failed to open database", "error", err, "path", dbPath)
        }
        defer db.Close()

        tlog.Info("Database opened successfully", "path", dbPath)

        // ... rest of test
    }

Available methods:
- Fatal(msg, ...keyvals) - logs error and fails test immediately
- Fatalf(msg, ...keyvals) - same as Fatal (for consistency)
- Error(msg, ...keyvals) - logs error and marks test as failed, but continues
- Errorf(msg, ...keyvals) - same as Error
- Info(msg, ...keyvals) - logs info message
- Log(msg, ...keyvals) - same as Info
- Logf(msg, ...keyvals) - same as Info
- Debug(msg, ...keyvals) - logs debug message
- Warn(msg, ...keyvals) - logs warning message

Output format:
    path/to/file.go:123: LEVEL: message key1=val1, key2=val2 (HH:MM:SS.mmm)
*/
