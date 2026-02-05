package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

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
	// Parse command-line arguments
	command := "up" // default command
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	// Validate command
	switch command {
	case "up", "down", "version", "drop", "reset", "drop-all":
		// Valid commands with no additional args
	case "force":
		if len(os.Args) < 3 {
			logging.Error("force command requires a version number")
			fmt.Println("Usage: migrate [up|down|version|drop|reset|drop-all|force <version>]")
			os.Exit(1)
		}
	default:
		logging.Error("Invalid command", "command", command)
		fmt.Println("Usage: migrate [up|down|version|drop|reset|drop-all|force <version>]")
		os.Exit(1)
	}

	logging.Info("Starting database migration...", "command", command)
	db, err := database.OpenDB()
	if err != nil {
		logging.Error("Error opening database", "error", err)
		os.Exit(1)
	}
	defer db.SQL.Close()

	driver, err := postgres.WithInstance(db.SQL, &postgres.Config{
		MigrationsTable: "schema_migrations",
		DatabaseName:    "kadeem",
	})
	if err != nil {
		logging.Error("Error creating database driver", "error", err)
		os.Exit(1)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres", driver)
	if err != nil {
		logging.Error("Error creating migration instance", "error", err)
		os.Exit(1)
	}

	m.Log = &migrateLogger{}

	// Execute the appropriate command
	switch command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			logging.Error("Error applying migrations", "error", err)
			os.Exit(1)
		}
		logging.Info("Database migrations applied successfully")

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			logging.Error("Error rolling back migration", "error", err)
			os.Exit(1)
		}
		logging.Info("Migration rolled back successfully")

	case "force":
		versionStr := os.Args[2]

		// Special case: force 0 or force -1 means reset migration tracking
		if versionStr == "0" || versionStr == "-1" {
			logging.Info("Version 0/-1 detected, resetting migration tracking...")
			_, err := db.SQL.Exec("DROP TABLE IF EXISTS schema_migrations")
			if err != nil {
				logging.Error("Error dropping schema_migrations table", "error", err)
				os.Exit(1)
			}
			logging.Info("Migration tracking reset successfully (use 'migrate up' to apply all migrations)")
			return
		}

		version, err := strconv.Atoi(versionStr)
		if err != nil {
			logging.Error("Invalid version number", "version", versionStr, "error", err)
			os.Exit(1)
		}
		if err := m.Force(version); err != nil {
			logging.Error("Error forcing version", "version", version, "error", err)
			os.Exit(1)
		}
		logging.Info("Database version forced successfully", "version", version)

	case "version":
		version, dirty, err := m.Version()
		if err != nil && err != migrate.ErrNilVersion {
			logging.Error("Error getting migration version", "error", err)
			os.Exit(1)
		}
		if err == migrate.ErrNilVersion {
			logging.Info("No migrations have been applied yet")
		} else {
			status := "clean"
			if dirty {
				status = "dirty"
			}
			logging.Info("Current migration version", "version", version, "status", status)
		}

	case "drop":
		if err := m.Drop(); err != nil {
			logging.Error("Error dropping database schema", "error", err)
			os.Exit(1)
		}
		logging.Info("Database schema dropped successfully")

	case "reset":
		// Just drop the schema_migrations table to reset tracking without destroying data
		_, err := db.SQL.Exec("DROP TABLE IF EXISTS schema_migrations")
		if err != nil {
			logging.Error("Error dropping schema_migrations table", "error", err)
			os.Exit(1)
		}
		logging.Info("Migration tracking reset successfully (schema_migrations table dropped)")

	case "drop-all":
		// Drop all tables in the public schema
		logging.Info("Dropping all tables in public schema...")
		_, err := db.SQL.Exec(`
			DO $$ DECLARE
				r RECORD;
			BEGIN
				FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
					EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
				END LOOP;
			END $$;
		`)
		if err != nil {
			logging.Error("Error dropping all tables", "error", err)
			os.Exit(1)
		}
		logging.Info("All tables dropped successfully")
	}
}

// Simple logger that implements migrate.Logger interface
type migrateLogger struct{}

func (l *migrateLogger) Printf(format string, v ...interface{}) {
	logging.Info(fmt.Sprintf("[MIGRATE] "+format, v...))
}

func (l *migrateLogger) Verbose() bool {
	return true // Enable verbose logging
}
