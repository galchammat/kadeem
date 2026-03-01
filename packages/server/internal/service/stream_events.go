package service

import (
	"fmt"

	"github.com/galchammat/kadeem/internal/database"
	"github.com/galchammat/kadeem/internal/model"
	"github.com/galchammat/kadeem/internal/twitch"
)

// StreamEventsService manages stream event syncing and retrieval.
type StreamEventsService struct {
	db     *database.DB
	twitch *twitch.TwitchClient
}

// NewStreamEventsService creates a new StreamEventsService.
func NewStreamEventsService(db *database.DB, twitchClient *twitch.TwitchClient) *StreamEventsService {
	return &StreamEventsService{db: db, twitch: twitchClient}
}

// SyncChannelEvents fetches and persists hype train and clip events for the given channel.
func (s *StreamEventsService) SyncChannelEvents(channelID string) error {
	hypeEvents, err := s.twitch.FetchHypeTrainEvents(channelID)
	if err != nil {
		return fmt.Errorf("fetch hype train events for channel %s: %w", channelID, err)
	}

	clipEvents, err := s.twitch.FetchTopClips(channelID)
	if err != nil {
		return fmt.Errorf("fetch clip events for channel %s: %w", channelID, err)
	}

	all := append(hypeEvents, clipEvents...)
	if err := s.db.UpsertStreamEvents(all); err != nil {
		return fmt.Errorf("upsert stream events for channel %s: %w", channelID, err)
	}
	return nil
}

// ListChannelEvents returns stream events for a specific channel within the given time range.
func (s *StreamEventsService) ListChannelEvents(channelID string, from, to int64, limit, offset int) ([]model.StreamEvent, error) {
	filter := &model.StreamEventFilter{
		ChannelID:    &channelID,
		TimestampMin: &from,
		TimestampMax: &to,
	}
	return s.db.ListStreamEvents(filter, limit, offset)
}

// ListStreamerEvents returns stream events for all channels of a streamer within the given time range.
func (s *StreamEventsService) ListStreamerEvents(streamerID int64, from, to int64, limit, offset int) ([]model.StreamEvent, error) {
	filter := &model.StreamEventFilter{
		StreamerID:   &streamerID,
		TimestampMin: &from,
		TimestampMax: &to,
	}
	return s.db.ListStreamEvents(filter, limit, offset)
}
