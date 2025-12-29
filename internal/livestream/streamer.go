package livestream

import (
	"context"

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

func NewStreamerClient(ctx context.Context, db *database.DB) *StreamerClient {
	return &StreamerClient{
		db:      *db,
		youtube: youtube.NewYoutubeClient(ctx),
		twitch:  twitch.NewTwitchClient(ctx),
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
		var lastLive int64
		channels, err := c.db.ListChannels(&models.Channel{StreamerID: streamer.ID})
		if err != nil {
			logging.Error("Failed to list channels for streamer", "streamer_id", streamer.ID, "error", err)
			return nil, err
		}
		for _, channel := range channels {
			broadcasts, err := c.ListBroadcasts(&models.Broadcast{ChannelID: channel.ID}, 1, 0)
			if err != nil {
				logging.Error("Failed to get latest broadcast", "channel_id", channel.ID, "error", err)
				return nil, err
			}
			if (len(broadcasts) > 0) && broadcasts[0].StartedAt > lastLive {
				lastLive = broadcasts[0].StartedAt
			}
			streamerView.Channels = append(streamerView.Channels, channel)
		}
		streamerView.LastLive = &lastLive
		streamerViews = append(streamerViews, streamerView)
	}
	return streamerViews, nil
}

func (c *StreamerClient) AddStreamer(name string) (bool, error) {
	streamer := models.Streamer{
		Name: name,
	}
	saved, err := c.db.SaveStreamer(streamer)
	if err != nil {
		logging.Error("Failed to add streamer", "error", err)
		return saved, err
	}
	return saved, nil
}

func (c *StreamerClient) DeleteStreamer(name string) (bool, error) {
	deleted, err := c.db.DeleteStreamer(name)
	if err != nil {
		logging.Error("Failed to delete streamer", "error", err)
		return deleted, err
	}
	return deleted, nil
}
