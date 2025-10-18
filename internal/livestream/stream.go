package livestream

import (
	"fmt"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/models"
)

func (c *StreamerClient) AddChannel(streamInput models.Stream) error {
	streams, err := c.db.ListStreams(&streamInput)
	if err != nil {
		errMessage := fmt.Errorf("failed to add stream while checking if stream already exists: %w", err)
		logging.Error(errMessage.Error())
		return errMessage
	}
	if len(streams) > 0 {
		errMessage := fmt.Errorf("stream already exists: %+v", streams[0])
		logging.Error(errMessage.Error())
		return errMessage
	}

	var stream models.Stream
	switch streamInput.Platform {
	case "twitch":
		stream, err = c.twitch.FindStream(streamInput)
	case "youtube":
		stream, err = c.youtube.FindStream(streamInput)
	default:
		errMessage := fmt.Errorf("unsupported platform: %s", streamInput.Platform)
		logging.Error(errMessage.Error())
		return errMessage
	}
	if err != nil {
		errMessage := fmt.Errorf("failed to find stream: %w", err)
		logging.Error(errMessage.Error())
		return errMessage
	}

	err = c.db.SaveStream(stream)
	if err != nil {
		errMessage := fmt.Errorf("failed to save stream: %w", err)
		logging.Error(errMessage.Error())
		return errMessage
	}
	return nil
}
