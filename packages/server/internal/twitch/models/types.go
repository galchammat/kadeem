package models

import (
	"encoding/json"
	"time"
)

type Streamer struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type Channel struct {
	ID          string `json:"id" db:"id"`
	StreamerID  int64  `json:"streamerId" db:"streamer_id"`
	Platform    string `json:"platform" db:"platform"`
	ChannelName string `json:"channelName" db:"channel_name"`
	AvatarURL   string `json:"avatarUrl" db:"avatar_url"`
	SyncedAt    *int64 `json:"syncedAt" db:"synced_at"`
}

type ChannelFilter struct {
	ID          *string
	StreamerID  *int64
	Platform    *string
	ChannelName *string
}

type StreamerView struct {
	StreamerID   int64     `json:"id" db:"id"`
	StreamerName string    `json:"name" db:"name"`
	Channels     []Channel `json:"channels" db:"channels"`
	LastLive     *int64    `json:"lastLive,omitempty" db:"last_live"`
	AvatarURL    *string   `json:"avatarUrl,omitempty" db:"avatar_url"`
}

type ChannelSearchResponse []struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"thumbnail_url"`
}

type APIResponse struct {
	Data       json.RawMessage `json:"data"`
	Pagination *struct {
		Cursor string `json:"cursor"`
	} `json:"pagination,omitempty"`
}

type DurationSeconds int64

func (d *DurationSeconds) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	duration, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = DurationSeconds(duration.Seconds())
	return nil
}
