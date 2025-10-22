package tests

import (
	"context"
	"os"
	"testing"

	"github.com/galchammat/kadeem/internal/database"
	"github.com/galchammat/kadeem/internal/livestream"
	"github.com/galchammat/kadeem/internal/models"
	"github.com/stretchr/testify/assert"
)

func testAddStreamer(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test; set RUN_INTEGRATION_TESTS=true to run it")
	}
	db, err := database.OpenDB()
	if err != nil {
		t.Fatal("Failed to connect to test database:", err)
	}
	defer db.SQL.Close()
	saved, err := db.SaveStreamer(models.Streamer{
		Name: "tarzaned",
	})
	if err != nil {
		t.Fatal("Failed to add streamer:", err)
	}
	if saved {
		t.Log("Streamer 'tarzaned' added to database")
	} else {
		t.Log("Streamer 'tarzaned' already exists in database")
	}
}

func testAddTwitchChannel(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test; set RUN_INTEGRATION_TESTS=true to run it")
	}

	ctx := context.Background()
	db, err := database.OpenDB()
	if err != nil {
		t.Fatal("Failed to connect to test database:", err)
	}
	defer db.SQL.Close()
	client := livestream.NewStreamerClient(ctx, db)

	// get StreamerID
	streamers, err := db.ListStreamers()
	if err != nil {
		t.Fatal("Failed to list streamers:", err)
	}
	var streamerID int64
	for _, s := range streamers {
		if s.Name == "tarzaned" {
			streamerID = s.ID
			break
		}
	}
	if streamerID == 0 {
		t.Fatal("Test streamer 'tarzaned' not found in database")
	}

	testChannelInput := models.Channel{
		StreamerID:  streamerID,
		Platform:    "twitch",
		ChannelName: "tarzaned",
	}
	saved, err := client.AddChannel(testChannelInput)
	assert.NoError(t, err, "Failed to add Twitch account")
	if saved {
		t.Log("Twitch account 'tarzaned' added to database")
	} else {
		t.Log("Twitch account 'tarzaned' already exists in database")
	}
}

func TestStreamChannel(t *testing.T) {
	t.Run("AddStreamer", testAddStreamer)
	t.Run("AddTwitchChannelToStreamer", testAddTwitchChannel)
}
