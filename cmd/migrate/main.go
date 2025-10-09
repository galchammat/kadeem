package main

import (
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"

	"github.com/galchammat/kadeem/internal/database"
	"github.com/galchammat/kadeem/internal/logging"
)

func init() {
	err := godotenv.Load(os.Getenv("ENV_FILE"))
	if err != nil {
		logging.Error("Error loading .env file", "error", err)
		os.Exit(1)
	}
}

func main() {
	logging.Info("Starting database migration...")
	db, err := database.OpenDB()
	if err != nil {
		logging.Error("Error opening database", "error", err)
		os.Exit(1)
	}
	defer db.SQL.Close()

	driver, err := sqlite3.WithInstance(db.SQL, &sqlite3.Config{
		MigrationsTable: "schema_migrations",
		DatabaseName:    "kadeem",
	})
	if err != nil {
		logging.Error("Error creating database driver", "error", err)
		os.Exit(1)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"sqlite3", driver)
	if err != nil {
		logging.Error("Error creating migration instance", "error", err)
		os.Exit(1)
	}

	m.Log = &migrateLogger{}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logging.Error("Error applying migrations", "error", err)
		os.Exit(1)
	}
	logging.Info("Database migrations applied successfully")
}

// Simple logger that implements migrate.Logger interface
type migrateLogger struct{}

func (l *migrateLogger) Printf(format string, v ...interface{}) {
	logging.Info(fmt.Sprintf("[MIGRATE] "+format, v...))
}

func (l *migrateLogger) Verbose() bool {
	return true // Enable verbose logging
}
