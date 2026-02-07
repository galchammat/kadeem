package tests

import (
	"testing"

	"github.com/galchammat/kadeem/internal/database"
	"github.com/galchammat/kadeem/internal/model"
	riot "github.com/galchammat/kadeem/internal/riot/api"
	"github.com/galchammat/kadeem/internal/service"
	"github.com/galchammat/kadeem/tests/logs"
)

func TestListLolMatches(t *testing.T) {
	tlog := testlog.New(t)

	db, err := database.OpenDB()
	if err != nil {
		tlog.Fatalf("Failed to open database", "error", err)
	}
	defer db.SQL.Close()

	matchSvc := service.NewMatchService(db, riot.NewClient())
	testPuuid := "OXR0AfpBu2Z-fFGu8KCE1sNzJLJbTpgClA42okBn-VsEVTwjJwMZu306s5JTLBmxPkVe2SSBIGe9ww"

	account, err := db.GetRiotAccount(testPuuid)
	if err != nil {
		tlog.Fatalf("Failed to get riot account", "error", err)
	}

	filter := model.LolMatchFilter{PUUID: &testPuuid}
	matches, err := matchSvc.ListMatches(&filter, account, 10, 0)
	if err != nil {
		tlog.Fatalf("Error fetching matches", "error", err)
	}

	if len(matches) == 0 {
		tlog.Fatalf("Expected to fetch at least one match, got 0")
	}

	tlog.Logf("Fetched matches successfully", "count", len(matches))
}
