package twitch

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/model"
)

func (c *TwitchClient) FindChannel(streamInput model.Channel) (model.Channel, error) {
	// prefer a human-friendly name for search, fall back to ChannelID if provided
	query := strings.TrimSpace(streamInput.ChannelName)
	if query == "" {
		query = strings.TrimSpace(streamInput.ID)
	}
	if query == "" {
		return model.Channel{}, fmt.Errorf("missing channel query (channelName or channelID required)")
	}

	endpoint := fmt.Sprintf("/search/channels?query=%s", url.QueryEscape(query))
	data, statusCode, err := c.makeRequest(endpoint)
	if err != nil {
		// propagate the underlying error, include status for easier debugging
		return model.Channel{}, fmt.Errorf("twitch search request failed (status %d): %w", statusCode, err)
	}

	var ChannelSearchResult model.ChannelSearchResponse
	if err := json.Unmarshal(data.Data, &ChannelSearchResult); err != nil {
		logging.Error("Failed to unmarshal Twitch search response", "query", query, "error", err)
		return model.Channel{}, fmt.Errorf("invalid twitch response: %w", err)
	}

	if len(ChannelSearchResult) == 0 {
		return model.Channel{}, fmt.Errorf("no twitch channel named %q was found", query)
	}

	for _, ch := range ChannelSearchResult {
		if strings.EqualFold(ch.DisplayName, query) {
			result := model.Channel{
				StreamerID:  streamInput.StreamerID,
				Platform:    "twitch",
				ChannelName: ch.DisplayName,
				ID:          ch.ID,
				AvatarURL:   ch.AvatarURL,
			}
			return result, nil
		}

	}

	return model.Channel{}, fmt.Errorf("no exact match found for channel %q", query)
}
