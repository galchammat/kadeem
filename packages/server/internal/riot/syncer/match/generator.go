package matchsync

import (
	"context"
	"fmt"
	"time"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/riot/api"
	"github.com/galchammat/kadeem/internal/riot/models"
)

const matchIDPageSize = 49
const defaultLookbackDays = 1
const defaultLookback = 60 * 60 * 24 * defaultLookbackDays

type MatchStore interface {
	ListRiotAccounts(filter *models.Account, limit, offset int) ([]models.Account, error)
	UpdateRiotAccount(puuid string, updates map[string]any) (bool, error)
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

func (s *MatchSyncer) Sync(ctx context.Context) error {
	for offset := 0; ; offset += matchIDPageSize {
		accounts, err := s.store.ListRiotAccounts(nil, matchIDPageSize, offset)
		if err != nil {
			return fmt.Errorf("list riot accounts: %w", err)
		}
		if len(accounts) == 0 {
			return nil
		}

		for _, account := range accounts {
			if err := s.syncAccount(ctx, account); err != nil {
				return err
			}
		}
	}
}

func (s *MatchSyncer) syncAccount(ctx context.Context, account models.Account) error {
	startTime := account.SyncedAt
	if startTime == nil {
		defaultStartTime := time.Now().Unix() - defaultLookback
		startTime = &defaultStartTime
	}
	for start := 0; ; start += matchIDPageSize {
		if err := ctx.Err(); err != nil {
			return err
		}

		matchIDs, err := s.client.FetchMatchIDPage(account.PUUID, account.Region, startTime, start, matchIDPageSize)
		if err != nil {
			return fmt.Errorf("fetch match id page for puuid %q start %d: %w", account.PUUID, start, err)
		}

		if len(matchIDs) == 0 {
			return nil
		}

		// details and timelines
		s.processMatches(ctx, matchIDs)

		logging.Info("synced riot match id page", "puuid", account.PUUID, "start", start, "count", len(matchIDs))
	}
}
