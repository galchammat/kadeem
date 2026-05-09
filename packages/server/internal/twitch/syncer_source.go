package twitch

import (
	"context"
	"fmt"

	"github.com/galchammat/kadeem/internal/syncer"
	twitchapi "github.com/galchammat/kadeem/internal/twitch/api"
	"github.com/galchammat/kadeem/internal/twitch/models"
)

var _ syncer.Source[models.StreamEvent] = (*StreamEventSyncer)(nil)

type StreamEventStore interface {
	ListChannels(filter *models.ChannelFilter, limit, offset int) ([]models.Channel, error)
	UpsertStreamEvents(events []models.StreamEvent) error
}

type StreamEventSyncer struct {
	client *twitchapi.TwitchClient
	store  StreamEventStore
}

func NewStreamEventSyncer(client *twitchapi.TwitchClient, store StreamEventStore) (*StreamEventSyncer, error) {
	if client == nil {
		return nil, fmt.Errorf("twitch client is nil")
	}
	if store == nil {
		return nil, fmt.Errorf("stream event store is nil")
	}

	return &StreamEventSyncer{client: client, store: store}, nil
}

func (s *StreamEventSyncer) Sync(ctx context.Context) error {
	platform := "twitch"
	channels, err := s.store.ListChannels(&models.ChannelFilter{Platform: &platform}, 1000, 0)
	if err != nil {
		return fmt.Errorf("list twitch channels: %w", err)
	}

	for _, channel := range channels {
		if err := ctx.Err(); err != nil {
			return err
		}

		hypeEvents, err := s.client.FetchHypeTrainEvents(channel.ID)
		if err != nil {
			return fmt.Errorf("fetch hype train events for channel %q: %w", channel.ID, err)
		}

		clipEvents, err := s.client.FetchTopClips(channel.ID)
		if err != nil {
			return fmt.Errorf("fetch clip events for channel %q: %w", channel.ID, err)
		}

		events := append(hypeEvents, clipEvents...)
		if err := s.store.UpsertStreamEvents(events); err != nil {
			return fmt.Errorf("upsert stream events for channel %q: %w", channel.ID, err)
		}
	}

	return nil
}
