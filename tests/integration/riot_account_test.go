package tests

import (
	"context"
	"os"
	"testing"

	"github.com/galchammat/kadeem/internal/database"
	"github.com/galchammat/kadeem/internal/riot/api"
)

func testListRiotAccounts(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test; set RUN_INTEGRATION_TESTS=true to run it")
	}

	ctx := context.Background()

	db, err := database.OpenDB()
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.SQL.Close()

	riotClient := riot.NewRiotClient(ctx, db)
	accounts, err := riotClient.ListAccounts(nil)
	if err != nil {
		t.Fatalf("Failed to list accounts: %v", err)
	}

	t.Log("Riot accounts listed successfully:", accounts)
}

func testAddRiotAccount(t *testing.T) {
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
	riotClient := riot.NewRiotClient(ctx, db)
	err = riotClient.AddAccount(region, gameName, tagLine, 0)
	if err != nil {
		t.Fatalf("Failed to add account: %v", err)
	}

	t.Log("Riot account saved successfully")
}

func TestRiotAccounts(t *testing.T) {
	t.Run("ListRiotAccounts", testListRiotAccounts)
	t.Run("AddRiotAccount", testAddRiotAccount)
}
