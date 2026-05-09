package tests

import (
	"context"
	"os"
	"testing"

	platformdb "github.com/galchammat/kadeem/internal/platform/database"
	"github.com/galchammat/kadeem/internal/service"
	twitchapi "github.com/galchammat/kadeem/internal/twitch/api"
	twitchstore "github.com/galchammat/kadeem/internal/twitch/store"
)

func TestSyncStreamEvents(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test; set RUN_INTEGRATION_TESTS=true to run it")
	}

	const channelID = "104410477"

	db, err := platformdb.OpenDB()
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer db.SQL.Close()
	store := twitchstore.New(db)

	twitchClient := twitchapi.NewTwitchClient(context.Background())
	svc := service.NewStreamEventsService(store, twitchClient)

	t.Run("SyncChannelEvents", func(t *testing.T) {
		if err := svc.SyncChannelEvents(channelID); err != nil {
			t.Fatalf("SyncChannelEvents: %v", err)
		}
		t.Log("sync completed without error")
	})

	t.Run("ListChannelEventsAfterSync", func(t *testing.T) {
		const (
			from   int64 = 0
			to     int64 = 9_999_999_999
			limit        = 100
			offset       = 0
		)
		events, err := svc.ListChannelEvents(channelID, from, to, limit, offset)
		if err != nil {
			t.Fatalf("ListChannelEvents: %v", err)
		}
		t.Logf("channel %s: %d event(s) in db", channelID, len(events))
		for _, ev := range events {
			t.Logf("  type=%-12s ts=%d title=%q", ev.EventType, ev.Timestamp, ev.Title)
		}
	})
}
