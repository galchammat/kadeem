package twitch

import (
	"fmt"

	"github.com/galchammat/kadeem/internal/models"
)

func (c *TwitchClient) FindStream(streamInput models.Stream) (models.Stream, error) {
	endpoint := fmt.Sprintf("/streams/%s", streamInput.ChannelID)
	data, statusCode, err := c.makeRequest(endpoint)
	if err != nil {
		return models.Stream{}, err
	}
	if statusCode == 404 {
		return models.Stream{}, nil
	}

}
