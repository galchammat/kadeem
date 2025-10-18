package tests

import (
	"context"
	"os"
	"testing"

	"github.com/galchammat/kadeem/internal/database"
	"github.com/galchammat/kadeem/internal/livestream"
	"github.com/galchammat/kadeem/internal/models"
)

func TestAddTwitchAccount(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test; set RUN_INTEGRATION_TESTS=true to run it")
	}

	ctx := context.Background()
	db, err := database.OpenDB()
	if err != nil {
		t.Fatal("Failed to connect to test database:", err)
	}
	defer db.SQL.Close()
	client := livestream.NewStreamerClient(ctx, *db)

	testStreamInput := models.Stream{
		Platform:    "twitch",
		ChannelName: "tarzaned",
	}

	err = client.AddChannel(testStreamInput)
	// err = client.AddTwitchAccount("test_user", "test_oauth_token")
	// assert.NoError(t, err, "Failed to add Twitch account")
	// return
}
