package twitch

import "github.com/galchammat/kadeem/internal/models"

func (c *TwitchClient) ListBroadcastsByStream(streamID models.StreamID, limit int) ([]models.Broadcast, error) {
	// Placeholder implementation
	return []models.Broadcast{}, nil
}
