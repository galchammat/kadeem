package livestream

import (
	"fmt"
	"time"

	"github.com/galchammat/kadeem/internal/models"
)

const syncRefreshInMinutes = 5

func (c *StreamerClient) SyncBroadcasts(channel models.Channel) error {
	var startTime int64
	if channel.SyncedAt != nil {
		startTime = *channel.SyncedAt
	} else {
		startTime = 0
	}

	var err error
	switch channel.Platform {
	case "twitch":
		{
			err = c.twitch.FetchBroadcasts(channel)
		}
	default:
		{
			err = fmt.Errorf("unsupported platform: %s", channel.Platform)
		}
	}
	if err == nil {
		c.db.UpdateChannel(channel.ID, map[string]interface{}{"synced_at": time.Now().Unix()})
	}
	return err
}

func (c *StreamerClient) ListBroadcasts(filters *models.Broadcast, limit int, offset int) ([]models.Broadcast, error) {
	if filters == nil || filters.ChannelID == "" {
		return []models.Broadcast{}, fmt.Errorf("channelID must be specified in filters")
	}

	channels, err := c.db.ListChannels(&models.Channel{ID: filters.ChannelID})
	if err != nil || len(channels) == 0 {
		return nil, err
	}
	channel := channels[0]

	// Check if channel needs sync (never synced or stale)
	if channel.SyncedAt == nil || (offset == 0 && time.Since(time.Unix(*channel.SyncedAt, 0)) > syncRefreshInMinutes*time.Minute) {
		c.SyncBroadcasts(channel)
	}

	broadcasts, err := c.db.ListBroadcasts(filters, limit, &offset)
	if err != nil {
		return nil, err
	}
	return broadcasts, nil
}
