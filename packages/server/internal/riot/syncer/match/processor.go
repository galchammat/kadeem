package matchsync

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/galchammat/kadeem/internal/models"
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

type Result struct {
	MatchID int64
	Op      Op

	MatchSummary riotmodels.MatchSummary
	Participants []riotmodels.MatchParticipantSummary
	Events       []any
}

func (s *MatchSyncer) processMatches(
	ctx context.Context,
	fullMatchIDs []string,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobs := make(chan Job, len(fullMatchIDs)*2)
	results := make(chan Result, len(fullMatchIDs)*2)
	fatalErr := make(chan error, 1)

	const workerCount = 10
	var wg sync.WaitGroup

	for range workerCount {
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
					select {
					case fatalErr <- err:
					default:
					}
					cancel()
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
		// jobs <- Job{FullMatchID: fullMatchID, Op: Timeline}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	count := len(fullMatchIDs)
	summaries := make([]riotmodels.MatchSummary, 0, count)
	participants := make([]riotmodels.MatchParticipantSummary, 0, count)
	events := make([]any, 0, count)

	for result := range results {
		fmt.Println(result)
		summaries = append(summaries, result.MatchSummary)
		participants = append(participants, result.Participants...)
		events = append(events, result.Events...)
	}

	select {
	case err := <-fatalErr:
		return err
	default:
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	if err := s.store.SaveMatchSummaryBatch(ctx, summaries); err != nil {
		return err
	}
	if err := s.store.SaveMatchParticipantBatch(ctx, participants); err != nil {
		return err
	}
	// s.store.SaveMatchEventBatch(ctx, events)

	return nil
}

func (s *MatchSyncer) processJob(ctx context.Context, job Job) (Result, error) {
	parts := strings.Split(job.FullMatchID, "_")
	if len(parts) != 2 {
		return Result{}, fmt.Errorf("invalid full match id %q", job.FullMatchID)
	}
	region, rawMatchID := parts[0], parts[1]
	matchID, err := strconv.ParseInt(rawMatchID, 10, 64)
	if err != nil {
		return Result{}, fmt.Errorf("failed to parse matchID %s. %w", rawMatchID, err)
	}

	result := Result{
		MatchID: matchID,
		Op:      job.Op,
	}

	switch job.Op {
	case Details:
		matchDetails, err := s.client.FetchMatchDetails(matchID, region)
		if err != nil {
			result.MatchSummary.ID = matchID
			result.MatchSummary.Region = region
			result.MatchSummary.Status = models.StatusRetry
			// ToDo - set to models.StatusDLQ if already Retry
		} else {
			result.MatchSummary, result.Participants = mapMatchDetails(*matchDetails)
			result.MatchSummary.Status = models.StatusDone
		}
	case Timeline:
		// Timeline processing not yet implemented
	default:
		return Result{}, fmt.Errorf("unknown op %q", job.Op)
	}

	return result, nil
}
