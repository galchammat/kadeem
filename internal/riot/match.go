package riot

import (
	"strings"
	"time"

	"github.com/galchammat/kadeem/internal/constants"
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

	_, err = c.db.UpdateRiotAccount(account.PUUID, map[string]interface{}{"synced_at": time.Now().Unix()})
	if err != nil {
		return err
	}
	return nil
}

// ListMatches retrieves League of Legends matches with optional filtering
func (c *RiotClient) ListMatches(filter *models.LolMatchFilter, limit int, offset int) ([]models.LeagueOfLegendsMatch, error) {
	if filter.LeagueOfLegendsAccount != nil &&
		(filter.LeagueOfLegendsAccount.SyncedAt == nil || time.Since(time.Unix(*filter.LeagueOfLegendsAccount.SyncedAt, 0)) > constants.SyncRefreshInMinutes*time.Minute) {
		err := c.SyncMatches(*filter.LeagueOfLegendsAccount)
		if err != nil {
			logging.Error("Error syncing matches for account", "PUUID", filter.LeagueOfLegendsAccount.PUUID, "Error", err)
		}
	}

	matches, err := c.db.ListLolMatches(filter, limit, offset)
	if err != nil {
		return nil, err
	}
	return matches, nil
}
