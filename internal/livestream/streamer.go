package livestream

import (
	"context"
	"fmt"

	"github.com/galchammat/kadeem/internal/database"
	"github.com/galchammat/kadeem/internal/livestream/twitch"
	"github.com/galchammat/kadeem/internal/livestream/youtube"
	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/models"
)

type StreamerClient struct {
	db      database.DB
	youtube *youtube.YoutubeClient
	twitch  *twitch.TwitchClient
}

func NewStreamerClient(ctx context.Context, db database.DB) *StreamerClient {
	return &StreamerClient{
		db:      db,
		youtube: youtube.NewYoutubeClient(ctx),
		twitch:  twitch.NewTwitchClient(ctx),
	}
}

func (c *StreamerClient) listBroadcastsByStream(stream models.Stream, limit int) ([]models.Broadcast, error) {
	switch stream.Platform {
	case "twitch":
		{
			return c.twitch.ListBroadcastsByStream(stream.ID, limit)
		}
	case "youtube":
		{
			return c.youtube.ListBroadcastsByStream(stream.ID, limit)
		}
	default:
		{
			return nil, fmt.Errorf("unsupported platform: %s", stream.Platform)
		}
	}
}

func (c *StreamerClient) ListStreamersWithDetails() ([]models.StreamerView, error) {
	var streamerViews []models.StreamerView

	streamers, err := c.db.ListStreamers()
	if err != nil {
		logging.Error("Failed to list streamers", "error", err)
		return nil, err
	}

	for _, streamer := range streamers {
		var streamerView = models.StreamerView{
			StreamerID:   streamer.ID,
			StreamerName: streamer.Name,
		}

		streams, err := c.db.ListStreamsByStreamer(streamer.ID)
		if err != nil {
			logging.Error("Failed to list streams for streamer", "streamer_id", streamer.ID, "error", err)
			return nil, err
		}

		for _, stream := range streams {
			broadcasts, err := ListBroadcastsByStream(stream.ID, 1)
			if err != nil {
				logging.Error("Failed to get latest broadcast", "stream_id", stream.ID, "error", err)
				return nil, err
			}
			var latestBroadcast models.Broadcast
			if len(broadcasts) == 0 {
				latestBroadcast = models.Broadcast{}
			} else {
				latestBroadcast = broadcasts[0]
			}

			streamerView.Streams = append(streamerView.Streams, models.Stream{
				ID:              stream.ID,
				Platform:        stream.Platform,
				ChannelName:     stream.ChannelName,
				ChannelID:       stream.ChannelID,
				LatestBroadcast: latestBroadcast,
			})
		}

		streamerViews = append(streamerViews, streamerView)
	}

	return streamerViews, nil
}
