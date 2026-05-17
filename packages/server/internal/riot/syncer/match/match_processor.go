package matchsync

import (
	"context"
	"sync"

	riotmodels "github.com/galchammat/kadeem/internal/riot/models"
)

type Op string

const (
	Details  Op = "details"
	Timeline Op = "timeline"
)

type Job struct {
	MatchID string
	Op      Op
}

type Status string

const (
	Done  Status = "done"
	Retry Status = "retry"
	DLQ   Status = "dlq"
)

type Result struct {
	MatchID string
	Op      Op
	Status  Status
	Err     error

	Payload any
}

func (s *MatchSyncer) processMatches(
	ctx context.Context,
	matchIDs []string,
) ([]Result, error) {
	jobs := make(chan Job, len(matchIDs)*2)
	results := make(chan Result, len(matchIDs)*2)
	const workerCount = 10
	var wg sync.WaitGroup

	return nil, nil
}
