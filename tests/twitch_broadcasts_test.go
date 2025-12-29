package tests

import (
	"context"
	"testing"

	"github.com/galchammat/kadeem/internal/database"
	"github.com/galchammat/kadeem/internal/livestream"
	"github.com/galchammat/kadeem/internal/models"
)

func TestListBroadcasts(t *testing.T) {
	const channelID string = "104410477" // test channel
	const limit, offset int = 0, 10

	ctx := context.Background()
	db, err := database.OpenDB()
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.SQL.Close()

	c := livestream.NewStreamerClient(ctx, db)
	broadcasts, err := c.ListBroadcasts(&models.Broadcast{ChannelID: channelID}, limit, offset)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	t.Log("Retrieved broadcasts:", broadcasts)
}
