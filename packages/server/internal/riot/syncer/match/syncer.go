package matchsync

import (
	"context"
	"fmt"

	"github.com/galchammat/kadeem/internal/riot/api"
	"github.com/galchammat/kadeem/internal/riot/models"
)

type MatchStore interface {
	// account sync
	ListRiotAccounts(filter *models.Account, limit, offset int) ([]models.Account, error)
	UpdateRiotAccount(puuid string, updates map[string]any) (bool, error)
	// match details
	SaveMatchSummaryBatch(context.Context, []models.MatchSummary) error
	SaveMatchParticipantBatch(context.Context, []models.MatchParticipantSummary) error
	// SaveMatchEventBatch(context.Context, []models.MatchEvent)
}

type MatchSyncer struct {
	client *api.Client
	store  MatchStore
}

func NewMatchSyncer(client *api.Client, store MatchStore) (*MatchSyncer, error) {
	if client == nil {
		return nil, fmt.Errorf("riot client is nil")
	}
	if store == nil {
		return nil, fmt.Errorf("match id store is nil")
	}

	return &MatchSyncer{client: client, store: store}, nil
}
