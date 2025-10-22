package twitch

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/models"
)

func (c *TwitchClient) FindChannel(streamInput models.Channel) (models.Channel, error) {
	// prefer a human-friendly name for search, fall back to ChannelID if provided
	query := strings.TrimSpace(streamInput.ChannelName)
	if query == "" {
		query = strings.TrimSpace(streamInput.ChannelID)
	}
	if query == "" {
		return models.Channel{}, fmt.Errorf("missing channel query (channelName or channelID required)")
	}

	endpoint := fmt.Sprintf("/search/channels?query=%s", url.QueryEscape(query))
	data, statusCode, err := c.makeRequest(endpoint)
	if err != nil {
		// propagate the underlying error, include status for easier debugging
		return models.Channel{}, fmt.Errorf("twitch search request failed (status %d): %w", statusCode, err)
	}

	var ChannelSearchResult models.ChannelSearchResponse
	if err := json.Unmarshal(data, &ChannelSearchResult); err != nil {
		logging.Error("failed to unmarshal twitch search response", err)
		return models.Channel{}, fmt.Errorf("invalid twitch response: %w", err)
	}

	if len(ChannelSearchResult.Data) == 0 {
		return models.Channel{}, fmt.Errorf("no channel found for query %q", query)
	}

	for _, ch := range ChannelSearchResult.Data {
		if strings.EqualFold(ch.DisplayName, query) {
			result := models.Channel{
				StreamerID:  streamInput.StreamerID,
				Platform:    "twitch",
				ChannelName: ch.DisplayName,
				ChannelID:   ch.ID,
				AvatarURL:   ch.AvatarURL,
			}
			return result, nil
		}

	}

	return models.Channel{}, fmt.Errorf("no exact match found for channel %q", query)
}
