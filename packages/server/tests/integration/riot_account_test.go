package tests

import (
	"os"
	"testing"

	platformdb "github.com/galchammat/kadeem/internal/platform/database"
	riotapi "github.com/galchammat/kadeem/internal/riot/api"
	riotpostgres "github.com/galchammat/kadeem/internal/riot/postgres"
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
	store := riotpostgres.New(db)

	accountSvc := service.NewAccountService(store, riotapi.NewClient())
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
	store := riotpostgres.New(db)

	accountSvc := service.NewAccountService(store, riotapi.NewClient())
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
