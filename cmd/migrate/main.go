package main

import (
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
	db, err := database.OpenDB()
	if err != nil {
		logging.Error("Error opening database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{
		MigrationsTable: "go_migrations",
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
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logging.Error("Error applying migrations", "error", err)
		os.Exit(1)
	}
	logging.Info("Database migrations applied successfully")
}
