package livestream

import (
	"fmt"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/models"
)

func (c *StreamerClient) AddChannel(channelInput models.Channel) (bool, error) {
	var err error
	var channel models.Channel
	switch channelInput.Platform {
	case "twitch":
		channel, err = c.twitch.FindChannel(channelInput)
	default:
		err := fmt.Errorf("unsupported platform: %s", channelInput.Platform)
		logging.Error("Unsupported platform requested", "platform", channelInput.Platform)
		return false, err
	}
	if err != nil {
		return false, fmt.Errorf("failed to find channel: %w", err)
	}

	saved, err := c.db.SaveChannel(channel)
	if err != nil {
		return false, fmt.Errorf("failed to save channel: %w", err)
	}
	return saved, nil
}

func (c *StreamerClient) DeleteChannel(channelID string) (bool, error) {
	deleted, err := c.db.DeleteChannel(channelID)
	if err != nil {
		return false, fmt.Errorf("failed to delete channel: %w", err)
	}
	return deleted, nil
}
