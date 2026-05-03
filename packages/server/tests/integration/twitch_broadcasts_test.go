package tests

import (
	"context"
	"os"
	"testing"

	"github.com/galchammat/kadeem/internal/model"
	platformdb "github.com/galchammat/kadeem/internal/platform/database"
	"github.com/galchammat/kadeem/internal/service"
	"github.com/galchammat/kadeem/internal/twitch"
	twitchstore "github.com/galchammat/kadeem/internal/twitch/store"
)

func TestListBroadcasts(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test; set RUN_INTEGRATION_TESTS=true to run it")
	}
	const channelID string = "104410477" // test channel
	const limit, offset int = 10, 0

	db, err := platformdb.OpenDB()
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.SQL.Close()
	store := twitchstore.New(db)

	twitchClient := twitch.NewTwitchClient(context.Background())
	streamerSvc := service.NewStreamerService(store, twitchClient)
	broadcasts, err := streamerSvc.ListBroadcasts(&model.Broadcast{ChannelID: channelID}, limit, offset)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	t.Log("Retrieved broadcasts:", "broadcasts", broadcasts)
}
