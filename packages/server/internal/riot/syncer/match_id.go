package syncer

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/riot/api"
	"github.com/galchammat/kadeem/internal/riot/models"
)

const matchIDPageSize = 100
const defaultLookback = 60 * 60 * 24 * 14

type MatchIDStore interface {
	ListRiotAccounts(filter *models.Account, limit, offset int) ([]models.Account, error)
	SaveMatchIDs(ctx context.Context, matchIDs []int64) error
	UpdateRiotAccount(puuid string, updates map[string]any) (bool, error)
}

type MatchIDSyncer struct {
	client *api.Client
	store  MatchIDStore
}

func NewMatchIDSyncer(client *api.Client, store MatchIDStore) (*MatchIDSyncer, error) {
	if client == nil {
		return nil, fmt.Errorf("riot client is nil")
	}
	if store == nil {
		return nil, fmt.Errorf("match id store is nil")
	}

	return &MatchIDSyncer{client: client, store: store}, nil
}

func (s *MatchIDSyncer) Sync(ctx context.Context) error {
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

func (s *MatchIDSyncer) syncAccount(ctx context.Context, account models.Account) error {
	startTime := account.SyncedAt
	if startTime == nil {
		defaultStartTime := time.Now().Unix() - defaultLookback
		startTime = &defaultStartTime
	}
	for start := 0; ; start += matchIDPageSize {
		if err := ctx.Err(); err != nil {
			return err
		}

		page, err := s.client.FetchMatchIDPage(account.PUUID, account.Region, startTime, start, matchIDPageSize)
		if err != nil {
			return fmt.Errorf("fetch match id page for puuid %q start %d: %w", account.PUUID, start, err)
		}
		if len(page) == 0 {
			return nil
		}

		matchIDs, err := numericMatchIDs(page)
		if err != nil {
			return err
		}

		if err := s.store.SaveMatchIDs(ctx, matchIDs); err != nil {
			return fmt.Errorf("save match id page for puuid %q start %d: %w", account.PUUID, start, err)
		}

		if _, err := s.store.UpdateRiotAccount(account.PUUID, map[string]any{"synced_at": time.Now().Unix()}); err != nil {
			return fmt.Errorf("update riot account sync time for puuid %q: %w", account.PUUID, err)
		}

		logging.Info("synced riot match id page", "puuid", account.PUUID, "start", start, "count", len(matchIDs))
	}
}

func numericMatchIDs(matchIDs []string) ([]int64, error) {
	numericIDs := make([]int64, 0, len(matchIDs))
	for _, matchID := range matchIDs {
		parts := strings.Split(matchID, "_")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid riot match id %q", matchID)
		}

		numericID, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse riot match id %q: %w", matchID, err)
		}

		numericIDs = append(numericIDs, numericID)
	}

	return numericIDs, nil
}
