package riot

import (
	"context"
	"fmt"

	"github.com/galchammat/kadeem/internal/model"
	"github.com/galchammat/kadeem/internal/syncer"
)

var _ syncer.Source[model.LolMatch] = (*MatchSyncer)(nil)

type MatchStore interface {
	SaveMatches(ctx context.Context, matches []model.LolMatch) error
}

type MatchSyncer struct {
	client   *Client
	store    MatchStore
	accounts []model.LolAccount
}

func NewMatchSyncer(client *Client, store MatchStore, accounts []model.LolAccount) (*MatchSyncer, error) {
	if client == nil {
		return nil, fmt.Errorf("riot client is nil")
	}
	if store == nil {
		return nil, fmt.Errorf("match store is nil")
	}

	return &MatchSyncer{
		client:   client,
		store:    store,
		accounts: accounts,
	}, nil
}

func (s *MatchSyncer) Sync(ctx context.Context) error {
	var matches []model.LolMatch
	for _, account := range s.accounts {
		if err := ctx.Err(); err != nil {
			return err
		}

		matchIDs, err := s.client.FetchMatchIDs(account.PUUID, account.Region, account.SyncedAt)
		if err != nil {
			return fmt.Errorf("fetch match ids for puuid %q: %w", account.PUUID, err)
		}

		for _, matchID := range matchIDs {
			if err := ctx.Err(); err != nil {
				return err
			}

			detail, err := s.client.FetchMatchDetail(matchID, account.Region)
			if err != nil {
				return fmt.Errorf("fetch match detail %q: %w", matchID, err)
			}

			participants := detail.Info.Participants
			for i := range participants {
				participants[i].GameID = detail.Info.ID
			}

			matches = append(matches, model.LolMatch{
				Summary: model.LolMatchSummary{
					ID:        detail.Info.ID,
					StartedAt: &detail.Info.StartedAt,
					Duration:  &detail.Info.Duration,
				},
				Participants: participants,
			})
		}
	}

	return s.store.SaveMatches(ctx, matches)
}
