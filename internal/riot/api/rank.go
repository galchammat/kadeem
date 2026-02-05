package riot

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/models"
)

// RiotRankEntry represents the Riot API response for rank data
type RiotRankEntry struct {
	QueueType    string `json:"queueType"`
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	LeaguePoints int    `json:"leaguePoints"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
}

// Sync Rank fetches current rank for an account and stores it
func (c *RiotClient) SyncRank(account *models.LeagueOfLegendsAccount) error {
	// First, get summoner ID from PUUID
	endpoint := fmt.Sprintf("/lol/summoner/v4/summoners/by-puuid/%s", account.PUUID)
	url := c.buildURL(account.Region, endpoint)

	body, _, err := c.makeRequest(url)
	if err != nil {
		logging.Error("Failed to get summoner ID", "puuid", account.PUUID, "error", err)
		return fmt.Errorf("failed to get summoner ID: %w", err)
	}

	var summoner struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(body, &summoner); err != nil {
		return fmt.Errorf("failed to decode summoner response: %w", err)
	}

	// Fetch rank data
	rankEndpoint := fmt.Sprintf("/lol/league/v4/entries/by-summoner/%s", summoner.ID)
	rankURL := c.buildURL(account.Region, rankEndpoint)

	rankBody, _, err := c.makeRequest(rankURL)
	if err != nil {
		logging.Error("Failed to fetch rank", "summonerID", summoner.ID, "error", err)
		return fmt.Errorf("failed to fetch rank: %w", err)
	}

	var entries []RiotRankEntry
	if err := json.Unmarshal(rankBody, &entries); err != nil {
		return fmt.Errorf("failed to decode rank response: %w", err)
	}

	// Store rank snapshot for each queue type
	timestamp := time.Now().Unix()
	for _, entry := range entries {
		queueID := queueTypeToID(entry.QueueType)
		if queueID == 0 {
			continue // Skip unknown queue types
		}

		rank := &models.PlayerRank{
			PUUID:        account.PUUID,
			Timestamp:    timestamp,
			Tier:         entry.Tier,
			Rank:         entry.Rank,
			LeaguePoints: entry.LeaguePoints,
			Wins:         entry.Wins,
			Losses:       entry.Losses,
			QueueId:      queueID,
		}

		if err := c.db.InsertPlayerRank(rank); err != nil {
			logging.Error("Failed to insert rank", "puuid", account.PUUID, "queueID", queueID, "error", err)
			return fmt.Errorf("failed to insert rank: %w", err)
		}
	}

	logging.Debug("Synced rank for account", "puuid", account.PUUID, "numEntries", len(entries))
	return nil
}

func queueTypeToID(queueType string) int {
	switch queueType {
	case "RANKED_SOLO_5x5":
		return 420
	case "RANKED_FLEX_SR":
		return 440
	case "RANKED_FLEX_TT":
		return 470
	default:
		return 0
	}
}
