package tests

import (
	"context"
	"testing"

	"github.com/galchammat/kadeem/internal/database"
	"github.com/galchammat/kadeem/internal/models"
	"github.com/galchammat/kadeem/internal/riot"
	"github.com/galchammat/kadeem/internal/testlog"
)

func TestListLolMatches(t *testing.T) {
	tlog := testlog.New(t)
	ctx := context.Background()
	DB, err := database.OpenDB()
	if err != nil {
		tlog.Fatalf("Failed to open database", "error", err)
	}
	defer DB.SQL.Close()

	c := riot.NewRiotClient(ctx, DB)
	testPuuid := "OXR0AfpBu2Z-fFGu8KCE1sNzJLJbTpgClA42okBn-VsEVTwjJwMZu306s5JTLBmxPkVe2SSBIGe9ww"

	// Fetch account to enable syncing behavior
	account, err := DB.GetRiotAccount(testPuuid)
	if err != nil {
		tlog.Fatalf("Failed to get riot account", "error", err)
	}

	filter := models.LolMatchFilter{PUUID: &testPuuid}
	matches, err := c.ListMatches(&filter, account, 10, 0)

	if err != nil {
		tlog.Fatalf("Error fetching matches", "error", err)
	}

	if len(matches) == 0 {
		tlog.Fatalf("Expected to fetch at least one match, got 0")
	}

	tlog.Logf("Fetched matches successfully", "count", len(matches))
}
