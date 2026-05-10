package syncer

import (
	"context"
	"fmt"

	"github.com/galchammat/kadeem/internal/riot/api"
	"github.com/galchammat/kadeem/internal/riot/models"
)

type MatchDetailsStore interface {
	ClaimPendingMatchID(ctx context.Context) (id *int64, region string, err error)
	AckMatch(ctx context.Context, matchId int, region string) error
	NackMatch(ctx context.Context, matchId int, region string) error
	SaveMatchDetails(ctx context.Context, match models.Match) error
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

func (s *MatchDetailsSyncer) Sync(ctx context.Context) error {
	s.SyncMatch(ctx)
	return nil
}

func (s *MatchDetailsSyncer) SyncMatch(ctx context.Context) (bool, error) {
	matchID, region, err := s.store.ClaimPendingMatchID(ctx)
	if err != nil {
		return false, fmt.Errorf("claim pending match: %w", err)
	}
	if matchID == nil {
		return false, nil
	}
	matchResponse, err := s.client.FetchMatchDetails(*matchID, region)
	fmt.Println(matchResponse)
	return true, nil
}
