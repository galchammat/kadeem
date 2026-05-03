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

	"github.com/stretchr/testify/assert"
)

func testAddStreamer(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test; set RUN_INTEGRATION_TESTS=true to run it")
	}
	db, err := platformdb.OpenDB()
	if err != nil {
		t.Fatal("Failed to connect to test database:", err)
	}
	defer db.SQL.Close()
	store := twitchstore.New(db)

	id, err := store.SaveStreamer(model.Streamer{
		Name: "tarzaned",
	})
	if err != nil {
		t.Fatal("Failed to add streamer:", err)
	}
	if id > 0 {
		t.Log("Streamer 'tarzaned' added to database, id:", id)
	} else {
		t.Log("Streamer 'tarzaned' already exists in database")
	}
}

func testAddTwitchChannel(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test; set RUN_INTEGRATION_TESTS=true to run it")
	}

	db, err := platformdb.OpenDB()
	if err != nil {
		t.Fatal("Failed to connect to test database:", err)
	}
	defer db.SQL.Close()
	store := twitchstore.New(db)

	twitchClient := twitch.NewTwitchClient(context.Background())
	streamerSvc := service.NewStreamerService(store, twitchClient)

	streamers, err := store.ListStreamers(1000, 0)
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

	testChannelInput := model.Channel{
		StreamerID:  streamerID,
		Platform:    "twitch",
		ChannelName: "tarzaned",
	}
	saved, err := streamerSvc.AddChannel(testChannelInput)
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
