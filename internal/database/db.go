// ...existing code...
package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDB() (*sql.DB, error) {
	const (
		dsn        = "file:kadeem.db?cache=shared"
		busyMillis = 5000
	)

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed open sqlite: %w", err)
	}

	// For single-writer SQLite apps keep max open conns to 1.
	db.SetMaxOpenConns(1)

	// Apply sensible PRAGMAs. Use Exec so settings persist for the connection.
	if _, err := db.Exec(fmt.Sprintf("PRAGMA busy_timeout = %d;", busyMillis)); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed set busy_timeout: %w", err)
	}
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed enable foreign_keys: %w", err)
	}
	if _, err := db.Exec("PRAGMA key = 'tME-o:mncQYX$*S';"); err != nil {
		return nil, fmt.Errorf("failed set encryption key: %w", err)
	}

	return db, nil
}

// EnableWAL switches the DB to WAL journal mode and verifies it.
// Call this after OpenDB if you want WAL enabled.
func EnableWAL(db *sql.DB) error {
	if _, err := db.Exec("PRAGMA journal_mode = WAL;"); err != nil {
		return fmt.Errorf("failed set journal_mode WAL: %w", err)
	}
	var mode string
	if err := db.QueryRow("PRAGMA journal_mode;").Scan(&mode); err != nil {
		return fmt.Errorf("failed verify journal_mode: %w", err)
	}
	if strings.ToLower(mode) != "wal" {
		return fmt.Errorf("journal_mode not set to WAL, result: %q", mode)
	}
	return nil
}

// ...existing code...
