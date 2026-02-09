package service

import (
	"fmt"
	"time"

	"github.com/galchammat/kadeem/internal/constants"
	"github.com/galchammat/kadeem/internal/database"
	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/model"
	"github.com/galchammat/kadeem/internal/twitch"
)

type StreamerService struct {
	db     *database.DB
	twitch *twitch.TwitchClient
}

func NewStreamerService(db *database.DB, twitchClient *twitch.TwitchClient) *StreamerService {
	return &StreamerService{db: db, twitch: twitchClient}
}

func (s *StreamerService) ListStreamersWithDetails() ([]model.StreamerView, error) {
	var streamerViews []model.StreamerView
	streamers, err := s.db.ListStreamers()
	if err != nil {
		return nil, err
	}
	for _, streamer := range streamers {
		var streamerView = model.StreamerView{
			StreamerID:   streamer.ID,
			StreamerName: streamer.Name,
		}
		var lastLive int64
		channels, err := s.db.ListChannels(&model.ChannelFilter{StreamerID: &streamer.ID})
		if err != nil {
			return nil, err
		}
		for _, channel := range channels {
			broadcasts, err := s.ListBroadcasts(&model.Broadcast{ChannelID: channel.ID}, 1, 0)
			if err != nil {
				return nil, err
			}
			if (len(broadcasts) > 0) && broadcasts[0].CreatedAt > lastLive {
				lastLive = broadcasts[0].CreatedAt
			}
			streamerView.Channels = append(streamerView.Channels, channel)
		}
		streamerView.LastLive = &lastLive
		streamerViews = append(streamerViews, streamerView)
	}
	return streamerViews, nil
}

func (s *StreamerService) AddStreamer(name string) (int64, error) {
	streamer := model.Streamer{Name: name}
	return s.db.SaveStreamer(streamer)
}

func (s *StreamerService) DeleteStreamer(name string) (bool, error) {
	return s.db.DeleteStreamer(name)
}

func (s *StreamerService) AddChannel(channelInput model.Channel) (bool, error) {
	var channel model.Channel
	var err error
	switch channelInput.Platform {
	case "twitch":
		channel, err = s.twitch.FindChannel(channelInput)
	default:
		return false, fmt.Errorf("unsupported platform: %s", channelInput.Platform)
	}
	if err != nil {
		return false, fmt.Errorf("failed to find channel: %w", err)
	}
	return s.db.SaveChannel(channel)
}

func (s *StreamerService) DeleteChannel(channelID string) (bool, error) {
	return s.db.DeleteChannel(channelID)
}

func (s *StreamerService) SyncBroadcasts(channel model.Channel) error {
	var startTime int64
	if channel.SyncedAt != nil {
		startTime = *channel.SyncedAt
	}
	logging.Info("Syncing broadcasts for channel", "ID", channel.ID, "name", channel.ChannelName)

	var broadcasts []model.Broadcast
	var err error
	switch channel.Platform {
	case "twitch":
		broadcasts, err = s.twitch.FetchBroadcasts(channel.ID, startTime)
	default:
		return fmt.Errorf("unsupported platform: %s", channel.Platform)
	}
	if err != nil {
		return err
	}

	if err := s.db.InsertBroadcasts(broadcasts); err != nil {
		return err
	}

	_, err = s.db.UpdateChannel(channel.ID, map[string]any{"synced_at": time.Now()})
	return err
}

func (s *StreamerService) ListBroadcasts(filter *model.Broadcast, limit, offset int) ([]model.Broadcast, error) {
	if filter == nil || filter.ChannelID == "" {
		return nil, fmt.Errorf("channelID must be specified")
	}

	channels, err := s.db.ListChannels(&model.ChannelFilter{ID: &filter.ChannelID})
	if err != nil {
		return nil, err
	}
	if len(channels) == 0 {
		return nil, fmt.Errorf("channel not found: %s", filter.ChannelID)
	}
	channel := channels[0]

	// Auto-sync if never synced or stale (only on first page)
	if channel.SyncedAt == nil || (offset == 0 && time.Since(time.Unix(*channel.SyncedAt, 0)) > constants.SyncRefreshInMinutes*time.Minute) {
		if err := s.SyncBroadcasts(channel); err != nil {
			return nil, err
		}
	}

	return s.db.ListBroadcasts(filter, limit, offset)
}
