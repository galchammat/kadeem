package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/galchammat/kadeem/internal/logging"
	platformdb "github.com/galchammat/kadeem/internal/platform/database"
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
	client := twitchapi.NewTwitchClient(ctx)
	if err := syncTwitchEvents(ctx, client, store); err != nil {
		logging.Error("failed to sync twitch events", "error", err)
		os.Exit(1)
	}

	logging.Info("synced twitch events")
}

type twitchEventStore interface {
	ListChannels(filter *twitchmodels.ChannelFilter, limit, offset int) ([]twitchmodels.Channel, error)
	UpsertStreamEvents(events []twitchmodels.StreamEvent) error
}

func syncTwitchEvents(ctx context.Context, client *twitchapi.TwitchClient, store twitchEventStore) error {
	platform := "twitch"
	channels, err := store.ListChannels(&twitchmodels.ChannelFilter{Platform: &platform}, 1000, 0)
	if err != nil {
		return fmt.Errorf("list twitch channels: %w", err)
	}

	for _, channel := range channels {
		if err := ctx.Err(); err != nil {
			return err
		}

		hypeEvents, err := client.FetchHypeTrainEvents(channel.ID)
		if err != nil {
			return fmt.Errorf("fetch hype train events for channel %q: %w", channel.ID, err)
		}

		clipEvents, err := client.FetchTopClips(channel.ID)
		if err != nil {
			return fmt.Errorf("fetch clip events for channel %q: %w", channel.ID, err)
		}

		events := append(hypeEvents, clipEvents...)
		if err := store.UpsertStreamEvents(events); err != nil {
			return fmt.Errorf("upsert stream events for channel %q: %w", channel.ID, err)
		}
	}

	return nil
}
