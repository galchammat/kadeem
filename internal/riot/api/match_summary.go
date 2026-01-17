package riot

import (
	"encoding/json"
	"fmt"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/models"
)

type matchSummaryResponse struct {
	Info struct {
		ID           int64                                           `json:"gameId"`
		StartedAt    int64                                           `json:"gameStartTimestamp"`
		Duration     int                                             `json:"gameDuration"`
		Participants []models.LeagueOfLegendsMatchParticipantSummary `json:"participants"`
	} `json:"info"`
}

func (c *RiotClient) FetchMatchSummary(puuid string) ([]string, error) {
	if puuid == "" {
		return nil, fmt.Errorf("puuid cannot be empty")
	}

	account, err := c.db.GetRiotAccount(puuid)
	if err != nil {
		return nil, err
	}

	url := c.buildURL(account.Region, fmt.Sprintf("/lol/match/v5/matches/by-puuid/%s/ids", puuid))
	body, _, err := c.makeRequest(url)
	if err != nil {
		logging.Error("Failed to fetch match IDs from Riot API", "puuid", puuid, "url", url, "error", err)
		return nil, err
	}

	var matchIDs []string
	if err := json.Unmarshal(body, &matchIDs); err != nil {
		logging.Error("Failed to unmarshal match IDs", "error", err)
		return nil, err
	}

	return matchIDs, nil
}

func (c *RiotClient) SyncMatchSummary(matchID int64, fullMatchID string, region string) error {
	if matchID == 0 {
		return fmt.Errorf("matchID cannot be zero")
	}

	url := c.buildURL(region, fmt.Sprintf("/lol/match/v5/matches/%s", fullMatchID))
	body, _, err := c.makeRequest(url)
	if err != nil {
		logging.Error("Failed to fetch match summary from Riot API", "matchID", matchID, "fullMatchID", fullMatchID, "url", url, "error", err)
		return err
	}

	var response matchSummaryResponse
	if err := json.Unmarshal(body, &response); err != nil {
		logging.Error("Failed to unmarshal match summary", "matchID", matchID, "error", err)
		return err
	}

	summary := models.LeagueOfLegendsMatchSummary{
		ID:        response.Info.ID,
		StartedAt: &response.Info.StartedAt,
		Duration:  &response.Info.Duration,
	}

	// Insert match and all participants atomically in a transaction
	if err := c.db.InsertLolMatchWithParticipants(&summary, response.Info.Participants); err != nil {
		logging.Error(
			"Failed to insert match with participants (transaction rolled back)",
			"matchID", matchID,
			"fullMatchID", fullMatchID,
			"participantCount", len(response.Info.Participants),
			"error", err,
		)
		return err
	}

	logging.Debug(
		"Successfully synced match summary with participants",
		"matchID", matchID,
		"participantCount", len(response.Info.Participants),
	)

	return nil
}
