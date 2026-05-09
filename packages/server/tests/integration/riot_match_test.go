package tests

import (
	"os"
	"testing"

	platformdb "github.com/galchammat/kadeem/internal/platform/database"
	riotapi "github.com/galchammat/kadeem/internal/riot/api"
	riot "github.com/galchammat/kadeem/internal/riot/models"
	riotpostgres "github.com/galchammat/kadeem/internal/riot/postgres"
	"github.com/galchammat/kadeem/internal/service"
	"github.com/galchammat/kadeem/tests/logs"
)

func TestListLolMatches(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test; set RUN_INTEGRATION_TESTS=true to run it")
	}
	tlog := testlog.New(t)

	db, err := platformdb.OpenDB()
	if err != nil {
		tlog.Fatalf("Failed to open database", "error", err)
	}
	defer db.SQL.Close()
	store := riotpostgres.New(db)

	matchSvc := service.NewMatchService(store, riotapi.NewClient())
	testPuuid := "OXR0AfpBu2Z-fFGu8KCE1sNzJLJbTpgClA42okBn-VsEVTwjJwMZu306s5JTLBmxPkVe2SSBIGe9ww"

	account, err := store.GetRiotAccount(testPuuid)
	if err != nil {
		tlog.Fatalf("Failed to get riot account", "error", err)
	}

	filter := riot.MatchFilter{PUUID: &testPuuid}
	matches, err := matchSvc.ListMatches(&filter, account, 10, 0)
	if err != nil {
		tlog.Fatalf("Error fetching matches", "error", err)
	}

	if len(matches) == 0 {
		tlog.Fatalf("Expected to fetch at least one match, got 0")
	}

	tlog.Logf("Fetched matches successfully", "count", len(matches))
}
