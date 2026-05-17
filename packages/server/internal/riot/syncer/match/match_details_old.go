package matchsync

import (
	"context"
	"fmt"

	"github.com/galchammat/kadeem/internal/riot/api"
)

type MatchDetailsStore interface {
	ClaimPendingMatch(ctx context.Context) (id *int64, region *string, err error)
	AckMatch(ctx context.Context, matchId int64, region string) error
	NackMatch(ctx context.Context, matchId int64, region string) error
	SaveMatchDetails(ctx context.Context, matchID int64, region string, startedAt int64, duration int, queueID int) error
}

type MatchDetailsSyncer struct {
	client *api.Client
	store  MatchDetailsStore
}

func NewMatchDetailsSyncer(client *api.Client, store MatchDetailsStore) (*MatchDetailsSyncer, error) {
	if client == nil {
		return nil, fmt.Errorf("riot client is nil")
	}
	if store == nil {
		return nil, fmt.Errorf("match id store is nil")
	}

	return &MatchDetailsSyncer{client: client, store: store}, nil
}

func (s *MatchDetailsSyncer) WorkerLoop(ctx context.Context) error {
	for true {
		matchFound, err := s.SyncMatch(ctx)
		if err != nil {
			return fmt.Errorf("SyncMatch failed %w", err)
		}
		if !matchFound {
			return nil
		}
	}
	return nil
}

func (s *MatchDetailsSyncer) SyncMatch(ctx context.Context) (bool, error) {
	matchID, region, err := s.store.ClaimPendingMatch(ctx)
	if err != nil {
		return false, fmt.Errorf("claim pending match: %w", err)
	}
	if matchID == nil {
		return false, nil
	}

	matchResponse, err := s.client.FetchMatchDetails(*matchID, *region)
	if err != nil {
		s.store.NackMatch(ctx, *matchID, *region)
		return false, fmt.Errorf("FetchMatchDetails failed: %w", err)
	}
	fmt.Println(matchResponse)

	err = s.store.SaveMatchDetails(ctx, *matchID, *region, matchResponse.Info.StartedAt, matchResponse.Info.Duration, matchResponse.Info.QueueID)
	if err != nil {
		s.store.NackMatch(ctx, *matchID, *region)
		return false, fmt.Errorf("SaveMatchDetails failed: %w", err)
	}

	err = s.store.SaveMatchParticipants(ctx, *matchID, *region, matchResponse.Info.Participants)
	if err != nil {
		s.store.NackMatch(ctx, *matchID, *region)
		return false, fmt.Errorf("SaveMatchParticipants failed: %w", err)
	}

	s.store.AckMatch(ctx, *matchID, *region)
	return true, nil
}
