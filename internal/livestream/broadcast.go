package livestream

import (
	"fmt"
	"time"

	"github.com/galchammat/kadeem/internal/logging"
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
	logging.Debug("Syncing broadcasts for channel", "ID", channel.ID, "name", channel.ChannelName)
	var err error
	var broadcasts []models.Broadcast
	switch channel.Platform {
	case "twitch":
		{
			broadcasts, err = c.twitch.FetchBroadcasts(channel.ID, startTime)
		}
	default:
		{
			err = fmt.Errorf("unsupported platform: %s", channel.Platform)
		}
	}
	if err != nil {
		return err
	}
	err = c.db.InsertBroadcasts(broadcasts)
	if err != nil {
		return err
	}
	// _, err = c.db.UpdateChannel(channel.ID, map[string]interface{}{"synced_at": time.Now().Unix()})
	// if err != nil {
	// 	return err
	// }
	return nil
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
	logging.Debug("Fetching broadcasts for channel", "ID", channel.ID, "name", channel.ChannelName)

	// Check if channel needs sync (never synced or stale)
	if channel.SyncedAt == nil || (offset == 0 && time.Since(time.Unix(*channel.SyncedAt, 0)) > syncRefreshInMinutes*time.Minute) {
		err = c.SyncBroadcasts(channel)
		if err != nil {
			return []models.Broadcast{}, err
		}
	}

	broadcasts, err := c.db.ListBroadcasts(filters, limit, &offset)
	if err != nil {
		return nil, err
	}
	return broadcasts, nil
}
