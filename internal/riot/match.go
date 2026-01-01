package riot

import (
	"fmt"
	"strings"
	"time"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/models"
)

func (c *RiotClient) SyncMatches(account models.LeagueOfLegendsAccount) error {
	logging.Debug("Syncing matches for account", "ID", account.PUUID)

	// Fetch the latest replays and extract MatchIDs string slice from them.
	var err error
	var replayURLs []string
	replayURLs, err = c.FetchReplayURLs(account.PUUID, account.Region)
	if err != nil {
		return err
	}

	for _, url := range replayURLs {
		start := strings.LastIndex(url, "/") + 1
		end := strings.Index(url[start:], ".")
		matchID := url[start : start+end]
		existingMatch, err := c.db.GetLolMatch(matchID)
		if err != nil {
			logging.Error("Error checking existing match", "MatchID", matchID, "Error", err)
		}

		// Download the replay if (matchID record does not exist) or (row.replay==nil)
		if existingMatch == nil || existingMatch.ReplaySynced == nil {
			logging.Debug("Downloading replay", "MatchID", matchID, "URL", url)
			err = c.SyncMatchReplay(matchID, url)
			if err != nil {
				logging.Error("Error downloading replay", "MatchID", matchID, "Error", err)
			}
		}

		// Fetch the match summary if (matchID record does not exist) or (row.gameStartTimestamp==nil)
		if existingMatch == nil || existingMatch.StartedAt == nil {
			logging.Debug("Fetching match summary", "MatchID", matchID)
			err = c.SyncMatchSummary(matchID, account.Region)
			if err != nil {
				logging.Error("Error fetching match summary", "MatchID", matchID, "Error", err)
			}
		}
	}

	err = c.db.UpdateRiotAccount(account.PUUID, map[string]interface{}{"synced_at": time.Now().Unix()})
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
