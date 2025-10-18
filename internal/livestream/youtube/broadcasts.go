package youtube

import "github.com/galchammat/kadeem/internal/models"

func (c *YoutubeClient) ListBroadcastsByStream(streamID models.StreamID, limit int) ([]models.Broadcast, error) {
	// Placeholder implementation
	return []models.Broadcast{}, nil
}
