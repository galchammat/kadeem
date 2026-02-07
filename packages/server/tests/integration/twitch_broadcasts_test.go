package tests

import (
	"context"
	"testing"

	"github.com/galchammat/kadeem/internal/database"
	"github.com/galchammat/kadeem/internal/model"
	"github.com/galchammat/kadeem/internal/service"
	"github.com/galchammat/kadeem/internal/twitch"
)

func TestListBroadcasts(t *testing.T) {
	const channelID string = "104410477" // test channel
	const limit, offset int = 10, 0

	db, err := database.OpenDB()
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.SQL.Close()

	twitchClient := twitch.NewTwitchClient(context.Background())
	streamerSvc := service.NewStreamerService(db, twitchClient)
	broadcasts, err := streamerSvc.ListBroadcasts(&model.Broadcast{ChannelID: channelID}, limit, offset)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	t.Log("Retrieved broadcasts:", "broadcasts", broadcasts)
}
