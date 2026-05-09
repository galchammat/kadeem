package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/galchammat/kadeem/internal/logging"
	platformdb "github.com/galchammat/kadeem/internal/platform/database"
	"github.com/galchammat/kadeem/internal/syncer"
	"github.com/galchammat/kadeem/internal/twitch"
	twitchapi "github.com/galchammat/kadeem/internal/twitch/api"
	twitchmodels "github.com/galchammat/kadeem/internal/twitch/models"
	twitchstore "github.com/galchammat/kadeem/internal/twitch/store"
	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	logging.Init(os.Stderr, slog.LevelInfo)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := platformdb.OpenDB()
	if err != nil {
		logging.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.SQL.Close()

	store := twitchstore.New(db)
	source, err := twitch.NewStreamEventSyncer(twitchapi.NewTwitchClient(ctx), store)
	if err != nil {
		logging.Error("failed to create twitch event source", "error", err)
		os.Exit(1)
	}

	metadataSyncer, err := syncer.NewMetadataSyncer[twitchmodels.StreamEvent](syncer.MetadataSyncerConfig{Logger: slog.Default()}, source)
	if err != nil {
		logging.Error("failed to create metadata syncer", "error", err)
		os.Exit(1)
	}

	if err := metadataSyncer.RunOnce(ctx); err != nil {
		logging.Error("failed to sync twitch events", "error", err)
		os.Exit(1)
	}

	logging.Info("synced twitch events")
}
