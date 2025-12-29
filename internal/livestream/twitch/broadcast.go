package twitch

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/galchammat/kadeem/internal/models"
)

func (c *TwitchClient) FetchBroadcasts(channelID string, startTime int64) ([]models.Broadcast, err) {
	// Build query params
	params := url.Values{}
	params.Set("user_id", channelID)
	params.Set("first", "100")
	params.Set("started_at", strconv.FormatInt(startTime, 10))

	// Construct full URL
	endpoint := "/videos?" + params.Encode()

	response, statusCode, err := c.makeRequest(endpoint)
	if err != nil || statusCode != 200 {
		return []models.Broadcast{}, fmt.Errorf("failed to fetch broadcasts: status=%d, error=%w", statusCode, err)
	}

	var broadcasts models.BroadcastListResponse
	json.Unmarshal(response, &broadcasts)

	return broadcasts.Data, nil
}
