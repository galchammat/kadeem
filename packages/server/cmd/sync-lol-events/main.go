package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/platform/database"
	"github.com/galchammat/kadeem/internal/riot"
	"github.com/galchammat/kadeem/internal/riot/api"
	"github.com/galchammat/kadeem/internal/riot/models"
	"github.com/galchammat/kadeem/internal/riot/postgres"
	"github.com/galchammat/kadeem/internal/syncer"
	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	logging.Init(os.Stderr, slog.LevelInfo)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := database.OpenDB()
	if err != nil {
		logging.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.SQL.Close()

	store := postgres.New(db)
	accounts, err := store.ListRiotAccounts(nil, 1000, 0)
	if err != nil {
		logging.Error("failed to list accounts for sync", "error", err)
		os.Exit(1)
	}

	source, err := riot.NewMatchSyncer(api.NewClient(), store, accounts)
	if err != nil {
		logging.Error("failed to create lol event source", "error", err)
		os.Exit(1)
	}

	metadataSyncer, err := syncer.NewMetadataSyncer[models.Match](syncer.MetadataSyncerConfig{Logger: slog.Default()}, source)
	if err != nil {
		logging.Error("failed to create metadata syncer", "error", err)
		os.Exit(1)
	}

	if err := metadataSyncer.RunOnce(ctx); err != nil {
		logging.Error("failed to sync lol events", "error", err)
		os.Exit(1)
	}

	logging.Info("synced lol events", "accounts", len(accounts))
}
