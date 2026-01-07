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

	apiRegion, err := GetAPIRegion(account.Region)
	if err != nil {
		return nil, err
	}

	url := c.buildURL(apiRegion, fmt.Sprintf("/lol/match/v5/matches/by-puuid/%s/ids", puuid))
	body, _, err := c.makeRequest(url)
	if err != nil {
		return nil, err
	}

	var matchIDs []string
	if err := json.Unmarshal(body, &matchIDs); err != nil {
		logging.Error("Failed to unmarshal match IDs", "error", err)
		return nil, err
	}

	return matchIDs, nil
}

func (c *RiotClient) SyncMatchSummary(matchID string, region string) error {
	if matchID == "" {
		return fmt.Errorf("matchID cannot be empty")
	}

	url := c.buildURL(region, fmt.Sprintf("/lol/match/v5/matches/%s", matchID))
	body, _, err := c.makeRequest(url)
	if err != nil {
		return err
	}

	var response matchSummaryResponse
	if err := json.Unmarshal(body, &response); err != nil {
		logging.Error("Failed to unmarshal match summary", "error", err)
		return err
	}

	summary := models.LeagueOfLegendsMatchSummary{
		ID:        response.Info.ID,
		StartedAt: &response.Info.StartedAt,
		Duration:  &response.Info.Duration,
	}

	if err := c.db.InsertLolMatchSummary(&summary); err != nil {
		logging.Error("Failed to insert match summary", "matchID", matchID, "error", err)
		return err
	}

	for _, participant := range response.Info.Participants {
		if err := c.db.InsertLolMatchParticipantSummary(&participant); err != nil {
			logging.Error("Failed to insert match participant", "matchID", matchID, "participantID", participant.ParticipantID, "error", err)
			return err
		}
	}

	return nil
}
