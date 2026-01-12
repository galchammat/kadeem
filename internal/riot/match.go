package riot

import (
	"fmt"
	"strconv"
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
		var matchIdString string = url[start : start+end]
		matchIdInt, err := strconv.ParseInt(matchIdString, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse matchID as int. matchID: %s", matchIdString)
		}

		limit := 1
		offset := 0
		existingMatches, err := c.db.ListLolMatches(&models.LolMatchFilter{MatchID: &matchIdInt}, &limit, &offset)
		if err != nil {
			return fmt.Errorf("Error while checking for an existing match. MatchID: %d. Error: %w", matchIdInt, err)
		}
		var existingMatch *models.LeagueOfLegendsMatch
		if len(existingMatches) != 0 {
			existingMatch = &existingMatches[0]
		}

		// Download the replay if (matchID record does not exist) or (row.replay==nil)
		if existingMatch == nil || existingMatch.ReplayURL == nil {
			logging.Debug("Downloading replay", "MatchID", matchIdString, "URL", url)
			err = c.SyncMatchReplay(matchIdString, url)
			if err != nil {
				logging.Error("Error downloading replay", "MatchID", matchIdString, "Error", err)
			}
		}

		// Fetch the match summary if (matchID record does not exist) or (row.gameStartTimestamp==nil)
		if existingMatch == nil || existingMatch.Summary.StartedAt == nil {
			logging.Debug("Fetching match summary", "MatchID", matchIdString)
			err = c.SyncMatchSummary(matchIdString, account.Region)
			if err != nil {
				logging.Error("Error fetching match summary", "MatchID", matchIdString, "Error", err)
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
// If account is provided and needs syncing, it will sync matches before querying
func (c *RiotClient) ListMatches(filter *models.LolMatchFilter, account *models.LeagueOfLegendsAccount, limit int, offset int) ([]models.LeagueOfLegendsMatch, error) {
	// Check if account needs syncing
	if account != nil &&
		(account.SyncedAt == nil || time.Since(time.Unix(*account.SyncedAt, 0)) > constants.SyncRefreshInMinutes*time.Minute) {
		err := c.SyncMatches(*account)
		if err != nil {
			logging.Error("Error syncing matches for account", "PUUID", account.PUUID, "Error", err)
		}
	}

	matches, err := c.db.ListLolMatches(filter, &limit, &offset)
	if err != nil {
		return nil, err
	}
	return matches, nil
}
