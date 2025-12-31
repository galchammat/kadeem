package riot

import (
	"fmt"
	"time"

	"github.com/galchammat/kadeem/internal/constants"
	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/models"
)

func (c *RiotClient) FetchMatches(puuid string, startTime int64) ([]models.LeagueOfLegendsMatchSummary, error) {
	// Implementation to fetch matches from Riot API using puuid and startTime
	return []models.LeagueOfLegendsMatchSummary{}, nil
}

func (c *RiotClient) SyncMatches(account models.LeagueOfLegendsAccount) error {
	logging.Debug("Syncing matches for account", "ID", account.PUUID)

	var startTime int64
	if account.SyncedAt != nil {
		startTime = *account.SyncedAt
	} else {
		startTime = constants.DefaultSyncWindowInSeconds

	}

	var err error
	var matches []models.LeagueOfLegendsMatchSummary
	matches, err = c.FetchMatches(account.PUUID, startTime)

	if err != nil {
		return err
	}
	err = c.db.InsertMatches(matches)
	if err != nil {
		return err
	}
	_, err = c.UpdateAccount(account.ID, map[string]interface{}{"synced_at": time.Now().Unix()})
	if err != nil {
		return err
	}
	return nil
}

func (c *RiotClient) ListMatches(filters *models.LeagueOfLegendsMatchSummary, limit int, offset int) ([]models.Broadcast, error) {
	if filters == nil || filters.AccountID == "" {
		return []models.Broadcast{}, fmt.Errorf("accountID must be specified in filters")
	}

	accounts, err := c.ListAccounts(&models.LeagueOfLegendsAccount{ID: filters.AccountID})
	if err != nil || len(accounts) == 0 {
		return nil, err
	}
	account := accounts[0]

	// Check if account needs sync (never synced or stale)
	if account.SyncedAt == nil || (offset == 0 && time.Since(time.Unix(*account.SyncedAt, 0)) > constants.syncRefreshInMinutes*time.Minute) {
		err = c.SyncMatches(account)
		if err != nil {
			return []models.Broadcast{}, err
		}
	}

	matches, err := c.db.ListMatches(filters, limit, &offset)
	if err != nil {
		return nil, err
	}
	return matches, nil
}
