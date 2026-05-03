package tests

import (
	"os"
	"testing"

	platformdb "github.com/galchammat/kadeem/internal/platform/database"
	riot "github.com/galchammat/kadeem/internal/riot/api"
	riotstore "github.com/galchammat/kadeem/internal/riot/store"
	"github.com/galchammat/kadeem/internal/service"
)

func testListRiotAccounts(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test; set RUN_INTEGRATION_TESTS=true to run it")
	}

	db, err := platformdb.OpenDB()
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.SQL.Close()
	store := riotstore.New(db)

	accountSvc := service.NewAccountService(store, riot.NewClient())
	accounts, err := accountSvc.ListAccounts(nil)
	if err != nil {
		t.Fatalf("Failed to list accounts: %v", err)
	}

	t.Log("Riot accounts listed successfully:", accounts)
}

func testAddRiotAccount(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test; set RUN_INTEGRATION_TESTS=true to run it")
	}

	db, err := platformdb.OpenDB()
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.SQL.Close()
	store := riotstore.New(db)

	accountSvc := service.NewAccountService(store, riot.NewClient())
	err = accountSvc.AddAccount("NA", "the thirsty rock", "NA1", 0)
	if err != nil {
		t.Fatalf("Failed to add account: %v", err)
	}

	t.Log("Riot account saved successfully")
}

func TestRiotAccounts(t *testing.T) {
	t.Run("ListRiotAccounts", testListRiotAccounts)
	t.Run("AddRiotAccount", testAddRiotAccount)
}
