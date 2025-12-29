package livestream

import (
	"fmt"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/models"
)

func (c *StreamerClient) AddChannel(channelInput models.Channel) (bool, error) {
	// channels, err := c.db.ListChannels(&channelInput)
	// if err != nil {
	// 	errMessage := fmt.Errorf("failed to add channel while checking if channel already exists: %w", err)
	// 	logging.Error(errMessage.Error())
	// 	return false, errMessage
	// }
	// if len(channels) > 0 {
	// 	return false, nil
	// }
	var err error
	var channel models.Channel
	switch channelInput.Platform {
	case "twitch":
		channel, err = c.twitch.FindChannel(channelInput)
	// case "youtube":
	// 	channel, err = c.youtube.FindChannel(channelInput)
	default:
		errMessage := fmt.Errorf("unsupported platform: %s", channelInput.Platform)
		logging.Error(errMessage.Error())
		return false, errMessage
	}
	if err != nil {
		errMessage := fmt.Errorf("failed to find channel: %w", err)
		logging.Error(errMessage.Error())
		return false, errMessage
	}

	logging.Debug("Found channel", "channel", channel)
	saved, err := c.db.SaveChannel(channel)
	if err != nil {
		errMessage := fmt.Errorf("failed to save channel: %w", err)
		logging.Error(errMessage.Error())
		return false, errMessage
	}
	return saved, nil
}

func (c *StreamerClient) DeleteChannel(channelID string) (bool, error) {
	deleted, err := c.db.DeleteChannel(channelID)
	if err != nil {
		errMessage := fmt.Errorf("failed to delete channel: %w", err)
		logging.Error(errMessage.Error())
		return false, errMessage
	}
	return deleted, nil
}
