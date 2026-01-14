package riot

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/galchammat/kadeem/internal/constants"
	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/models"
)

func extractMatchID(url string) (int64, error) {
	re := regexp.MustCompile(`([A-Z0-9]+_\d+)\.replay`)
	matches := re.FindStringSubmatch(url)

	if len(matches) < 2 {
		return 0, fmt.Errorf("no match ID found in URL")
	}

	// Extract numeric portion after underscore
	fullMatchID := matches[1] // e.g., "EUW1_7665669531"
	parts := regexp.MustCompile(`_(\d+)$`).FindStringSubmatch(fullMatchID)
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid match ID format: %s", fullMatchID)
	}

	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid match ID: %v", err)
	}

	return id, nil
}

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
		matchID, err := extractMatchID(url)
		if err != nil {
			return fmt.Errorf("failed to parse int matchID from replay URL: %s", url)
		}

		if err != nil {
			return fmt.Errorf("failed to extract full matchID from replay URL: %s", url)
		}

		var limit, offset int = 1, 0
		existingMatches, err := c.db.ListLolMatches(&models.LolMatchFilter{MatchID: &matchID}, &limit, &offset)
		if err != nil {
			return fmt.Errorf("Error while checking for an existing match. MatchID: %d. Error: %w", matchID, err)
		}
		var existingMatch *models.LeagueOfLegendsMatch
		if len(existingMatches) != 0 {
			existingMatch = &existingMatches[0]
		}

		// Fetch the match summary if (matchID record does not exist) or (row.gameStartTimestamp==nil)
		if existingMatch == nil || existingMatch.Summary.StartedAt == nil {
			logging.Debug("Fetching match summary", "MatchID", matchID)
			err = c.SyncMatchSummary(matchID, account.Region)
			if err != nil {
				logging.Error("Error fetching match summary", "MatchID", matchID, "Error", err)
			}
		}

		// Download the replay if (matchID record does not exist) or (row.replay==nil)
		if existingMatch == nil || existingMatch.ReplayURL == nil {
			logging.Debug("Downloading replay", "MatchID", matchID, "URL", url)
			err = c.SyncMatchReplay(matchID, url)
			if err != nil {
				logging.Error("Error downloading replay", "MatchID", matchID, "Error", err)
			}
		}

	}

	_, err = c.db.UpdateRiotAccount(account.PUUID, map[string]interface{}{"synced_at": time.Now().Unix()})
	if err != nil {
		return err
	}
	return nil
}

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
