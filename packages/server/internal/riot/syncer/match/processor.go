package matchsync

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	riotmodels "github.com/galchammat/kadeem/internal/riot/models"
)

type Op string

const (
	Details  Op = "details"
	Timeline Op = "timeline"
)

type Job struct {
	FullMatchID string
	Region      string
	Op          Op
}

type Status string

const (
	Done  Status = "done"
	Retry Status = "retry"
	DLQ   Status = "dlq"
)

type Result struct {
	MatchID int64
	Region  string
	Op      Op
	Status  Status
	Err     error

	Payload any
}

func (s *MatchSyncer) processMatches(
	ctx context.Context,
	fullMatchIDs []string,
) ([]Result, error) {
	jobs := make(chan Job, len(fullMatchIDs)*2)
	results := make(chan Result, len(fullMatchIDs)*2)

	const workerCount = 10
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for job := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
				}

				result, err := s.processJob(ctx, job)
				if err != nil {
					return
				}

				select {
				case results <- result:
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	for _, fullMatchID := range fullMatchIDs {
		jobs <- Job{FullMatchID: fullMatchID, Op: Details}
		jobs <- Job{FullMatchID: fullMatchID, Op: Timeline}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	out := make([]Result, 0, len(fullMatchIDs)*2)
	for result := range results {
		out = append(out, result)
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *MatchSyncer) processJob(ctx context.Context, job Job) (Result, error) {
	parts := strings.Split(job.FullMatchID, "_")
	if len(parts) != 2 {
		return Result{}, fmt.Errorf("invalid full match id %q", job.FullMatchID)
	}
	region, rawMatchID := parts[0], parts[1]
	matchID, err := strconv.ParseInt(rawMatchID, 10, 64)
	if err != nil {
		return Result{}, fmt.Errorf("Failed to parse matchID %s. %w", rawMatchID, err)
	}

	result := Result{
		MatchID: matchID,
		Region:  region,
		Op:      job.Op,
	}

	switch job.Op {
	case Details:
		result.Payload, result.Err = s.client.FetchMatchDetails(matchID, region)
	case Timeline:
		result.Payload, result.Err = nil, nil
	default:
		result.Err = fmt.Errorf("unknown op %q", job.Op)
	}

	if result.Err != nil {
		result.Status = Retry
		// ToDo - set to DLQ if already Retry
	} else {
		result.Status = Done
	}

	return result, nil
}
