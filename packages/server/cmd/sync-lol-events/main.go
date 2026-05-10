package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/platform/database"
	"github.com/galchammat/kadeem/internal/riot/api"
	"github.com/galchammat/kadeem/internal/riot/postgres"
	"github.com/galchammat/kadeem/internal/riot/syncer"
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
	apiClient := api.NewClient()

	matchIDSyncer, err := syncer.NewMatchIDSyncer(apiClient, store)
	if err != nil {
		logging.Error("failed to create lol match id syncer", "error", err)
		os.Exit(1)
	}
	if err := matchIDSyncer.Sync(ctx); err != nil {
		logging.Error("failed to sync lol match ids", "error", err)
		os.Exit(1)
	}

	matchDetailsSyncer, err := syncer.NewMatchDetailsSyncer(apiClient, store)
	if err != nil {
		logging.Error("failed to create lol match replay syncer", "error", err)
		os.Exit(1)
	}
	if err := matchDetailsSyncer.WorkerLoop(ctx); err != nil {
		logging.Error("failed to sync lol match details", "error", err)
		os.Exit(1)
	}

	logging.Info("synced lol match ids")
}
