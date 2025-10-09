package database

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const defaultKey = "tME-o:mncQYX$*S"

type DB struct {
	SQL *sql.DB
}

func OpenDB() (*DB, error) {
	const filePath = "/home/galchammat/code/personal/kadeem/kadeem.db"
	const DSN = "file://" + filePath + "?cache=shared"
	const (
		dsn        = DSN
		busyMillis = 5000
		dbFile     = filePath
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

	// Decide whether to KEY (open encrypted DB) or REKEY (encrypt new DB / change key).
	// If the DB file exists and is non-empty, attempt PRAGMA key = '...'
	// Otherwise use PRAGMA rekey = '...' to encrypt/initialize the DB.
	fi, err := os.Stat(dbFile)
	keyEsc := strings.ReplaceAll(defaultKey, "'", "''")
	if err == nil && fi.Size() > 0 {
		// existing file -> provide key to decrypt
		if _, err := db.Exec("PRAGMA key = '" + keyEsc + "';"); err != nil {
			db.Close()
			return nil, fmt.Errorf("failed set encryption key (PRAGMA key): %w", err)
		}
		// verify key works by running a harmless query
		var n int
		if err := db.QueryRow("SELECT count(*) FROM sqlite_master;").Scan(&n); err != nil {
			db.Close()
			return nil, fmt.Errorf("encryption key verification failed: %w", err)
		}
	} else {
		// new or empty DB -> encrypt it (or set key). PRAGMA rekey will encrypt when SQLCipher is present.
		if _, err := db.Exec("PRAGMA rekey = '" + keyEsc + "';"); err != nil {
			db.Close()
			return nil, fmt.Errorf("failed set encryption key (PRAGMA rekey): %w", err)
		}
		// verify by running a harmless query
		var n int
		if err := db.QueryRow("SELECT count(*) FROM sqlite_master;").Scan(&n); err != nil {
			db.Close()
			return nil, fmt.Errorf("encryption rekey verification failed: %w", err)
		}
	}

	return &DB{SQL: db}, nil
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
