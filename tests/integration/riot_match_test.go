package tests

import (
	"context"
	"github.com/galchammat/kadeem/internal/database"
	"github.com/galchammat/kadeem/internal/models"
	"github.com/galchammat/kadeem/internal/riot"
	"testing"
)

func TestListLolMatches(t *testing.T) {
	ctx := context.Background()
	DB, err := database.OpenDB()
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer DB.SQL.Close()

	c := riot.NewRiotClient(ctx, DB)
	testPuuid := "OXR0AfpBu2Z-fFGu8KCE1sNzJLJbTpgClA42okBn-VsEVTwjJwMZu306s5JTLBmxPkVe2SSBIGe9ww"
	filter := models.LolMatchFilter{PUUID: &testPuuid}
	matches, err := c.ListMatches(&filter, 10, 0)

	if err != nil {
		t.Fatalf("Error fetching matches: %v", err)
	}

	if len(matches) == 0 {
		t.Fatalf("Expected to fetch at least one match, got 0")
	}

	t.Logf("Fetched %v matches successfully", matches)
}
