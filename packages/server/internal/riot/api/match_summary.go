package riot

import (
	"encoding/json"
	"fmt"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/model"
)

// MatchDetailResponse represents the Riot API response for match details.
type MatchDetailResponse struct {
	Info struct {
		ID           int64                                          `json:"gameId"`
		StartedAt    int64                                          `json:"gameStartTimestamp"`
		Duration     int                                            `json:"gameDuration"`
		Participants []model.LolMatchParticipantSummary `json:"participants"`
	} `json:"info"`
}

// FetchMatchIDs fetches match IDs for a PUUID from the Riot API.
// startTime is optional (unix timestamp in milliseconds, exclusive lower bound).
// Always uses count=100 (maximum allowed by Riot API).
func (c *Client) FetchMatchIDs(puuid, region string, startTime *int64) ([]string, error) {
	if puuid == "" {
		return nil, fmt.Errorf("puuid cannot be empty")
	}

	endpoint := fmt.Sprintf("/lol/match/v5/matches/by-puuid/%s/ids", puuid)
	url := c.buildURL(region, endpoint)

	// Add query parameters if provided
	query := ""
	if startTime != nil {
		query = fmt.Sprintf("?start=%d", *startTime)
	}
	// Always append count=100 (maximum allowed by Riot API)
	if query == "" {
		query = "?count=100"
	} else {
		query = fmt.Sprintf("%s&count=100", query)
	}
	url = url + query

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

// FetchMatchDetail fetches full match detail for a given match ID.
func (c *Client) FetchMatchDetail(fullMatchID, region string) (*MatchDetailResponse, error) {
	url := c.buildURL(region, fmt.Sprintf("/lol/match/v5/matches/%s", fullMatchID))
	body, _, err := c.makeRequest(url)
	if err != nil {
		logging.Error("Failed to fetch match detail from Riot API", "fullMatchID", fullMatchID, "url", url, "error", err)
		return nil, err
	}

	var response MatchDetailResponse
	if err := json.Unmarshal(body, &response); err != nil {
		logging.Error("Failed to unmarshal match detail", "fullMatchID", fullMatchID, "error", err)
		return nil, err
	}

	return &response, nil
}

// FetchReplayURLs fetches replay download URLs for a PUUID.
func (c *Client) FetchReplayURLs(puuid, region string) ([]string, error) {
	endpoint := fmt.Sprintf("/lol/match/v5/matches/by-puuid/%s/replays", puuid)
	url := c.buildURL(region, endpoint)
	body, statusCode, err := c.makeRequest(url)
	if err != nil || statusCode != 200 {
		logging.Error("Failed to fetch replay URLs from Riot API", "puuid", puuid, "region", region, "statusCode", statusCode, "error", err)
		return nil, fmt.Errorf("error fetching replay URLs: %v Status Code: %d", err, statusCode)
	}

	var replays model.LolApiReplaysResponse
	if err := json.Unmarshal(body, &replays); err != nil {
		logging.Error("Failed to unmarshal replay URLs response", "puuid", puuid, "error", err)
		return nil, err
	}
	return replays.URLs, nil
}

// RankEntry represents the Riot API response for rank data.
type RankEntry struct {
	QueueType    string `json:"queueType"`
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	LeaguePoints int    `json:"leaguePoints"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
}

// FetchSummonerID fetches the summoner ID for a PUUID.
func (c *Client) FetchSummonerID(puuid, region string) (string, error) {
	url := c.buildURL(region, fmt.Sprintf("/lol/summoner/v4/summoners/by-puuid/%s", puuid))
	body, _, err := c.makeRequest(url)
	if err != nil {
		logging.Error("Failed to get summoner ID", "puuid", puuid, "error", err)
		return "", fmt.Errorf("failed to get summoner ID: %w", err)
	}

	var summoner struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(body, &summoner); err != nil {
		return "", fmt.Errorf("failed to decode summoner response: %w", err)
	}

	return summoner.ID, nil
}

// FetchRankEntries fetches rank entries for a summoner.
func (c *Client) FetchRankEntries(summonerID, region string) ([]RankEntry, error) {
	url := c.buildURL(region, fmt.Sprintf("/lol/league/v4/entries/by-summoner/%s", summonerID))
	body, _, err := c.makeRequest(url)
	if err != nil {
		logging.Error("Failed to fetch rank", "summonerID", summonerID, "error", err)
		return nil, fmt.Errorf("failed to fetch rank: %w", err)
	}

	var entries []RankEntry
	if err := json.Unmarshal(body, &entries); err != nil {
		return nil, fmt.Errorf("failed to decode rank response: %w", err)
	}

	return entries, nil
}
