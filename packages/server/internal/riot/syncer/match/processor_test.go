package matchsync

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	riotapi "github.com/galchammat/kadeem/internal/riot/api"
)

func TestProcessDetailsJob(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test; set RUN_INTEGRATION_TESTS=true to run it")
	}

	ctx := context.Background()
	s := &MatchSyncer{client: riotapi.NewClient()}

	matchDetails, err := s.processJob(ctx, Job{
		FullMatchID: "EUW1_7848931380",
		Region:      "EUW1",
		Op:          Details,
	})

	if err != nil {
		t.Errorf("failed to process details job: %v", err)
	}

	b, err := json.MarshalIndent(matchDetails, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))

}
