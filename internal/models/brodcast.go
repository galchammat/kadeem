package models

import (
	"encoding/json"
	"time"
)

func (b *Broadcast) UnmarshalJSON(data []byte) error {
	type Alias Broadcast
	aux := struct {
		ID             interface{} `json:"id"` // explicitly ignore
		ChannelID      string      `json:"channel_id"`
		UserID         string      `json:"user_id"`
		CreatedAtStr   string      `json:"created_at"`
		PublishedAtStr string      `json:"published_at"`
		*Alias
	}{
		Alias: (*Alias)(b),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	switch {
	case aux.ChannelID != "":
		b.ChannelID = aux.ChannelID
	case aux.UserID != "":
		b.ChannelID = aux.UserID
	}
	// Parse created_at and published_at if present
	if aux.CreatedAtStr != "" {
		t, err := time.Parse(time.RFC3339, aux.CreatedAtStr)
		if err == nil {
			b.CreatedAt = t.Unix()
		}
	}
	if aux.PublishedAtStr != "" {
		t, err := time.Parse(time.RFC3339, aux.PublishedAtStr)
		if err == nil {
			b.PublishedAt = t.Unix()
		}
	}
	return nil
}

type Broadcast struct {
	ID           int64           `json:"id" db:"id"`
	ChannelID    string          `json:"channel_id" db:"channel_id"`
	Title        string          `json:"title" db:"title"`
	URL          string          `json:"url" db:"url"`
	ThumbnailURL string          `json:"thumbnail_url" db:"thumbnail_url"`
	Viewable     string          `json:"viewable" db:"viewable"`
	CreatedAt    int64           `json:"created_at" db:"created_at"`
	PublishedAt  int64           `json:"published_at" db:"published_at"`
	Duration     DurationSeconds `json:"duration" db:"duration"`
}
