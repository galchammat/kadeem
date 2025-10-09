package tests

import (
	"context"
	"os"
	"testing"

	"github.com/galchammat/kadeem/internal/database"
	"github.com/galchammat/kadeem/internal/riot"
)

func TestSaveRiotAccount(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test; set RUN_INTEGRATION_TESTS=true to run it")
	}

	ctx := context.Background()

	db, err := database.OpenDB()
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.SQL.Close()

	region := "NA"
	gameName := "the thirsty rock"
	tagLine := "NA1"
	riotClient := riot.NewRiotClient(ctx)
	account, err := riotClient.GetAccount(region, gameName, tagLine)
	if err != nil {
		t.Fatalf("Failed to get account: %v", err)
	}
	t.Log("Fetched Riot account", account)
	account.Region = region

	if err := db.SaveRiotAccount(&account); err != nil {
		t.Fatalf("Failed to save Riot account: %v", err)
	}

	t.Log("Riot account saved successfully")
}
