package twitch

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/models"
)

func (c *TwitchClient) FetchBroadcasts(channelID string, startTime int64) ([]models.Broadcast, error) {
	// Build query params
	params := url.Values{}
	params.Set("user_id", channelID)
	params.Set("first", "100")
	params.Set("type", "archive")

	// Construct full URL
	endpoint := "/videos?" + params.Encode()

	response, statusCode, err := c.makeRequest(endpoint)
	if err != nil || statusCode != 200 {
		return []models.Broadcast{}, fmt.Errorf("failed to fetch broadcasts: status=%d, error=%w", statusCode, err)
	}

	var rawMessages []json.RawMessage
	if err := json.Unmarshal(response.Data, &rawMessages); err != nil {
		return nil, fmt.Errorf("failed to unmarshal broadcasts: %w", err)
	}

	broadcasts := make([]models.Broadcast, 0, len(rawMessages))
	for _, raw := range rawMessages {
		var b models.Broadcast
		if err := json.Unmarshal(raw, &b); err != nil {
			logging.Warn("Failed to unmarshal broadcast", "error", err)
			continue
		}
		broadcasts = append(broadcasts, b)
	}

	return broadcasts, nil
}
