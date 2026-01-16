package models

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

// ChannelFilter provides filtering options for listing channels
type ChannelFilter struct {
	ID         *string `db:"id" op:"="`
	StreamerID *int64  `db:"streamer_id" op:"="`
	Platform   *string `db:"platform" op:"="`

	// Text search (user provides wildcards like "%search%")
	ChannelName *string `db:"channel_name" op:"like"`
}

type StreamerView struct {
	StreamerID   int64     `json:"id" db:"id"`
	StreamerName string    `json:"name" db:"name"`
	Channels     []Channel `json:"channels" db:"channels"`
	LastLive     *int64    `json:"lastLive,omitempty" db:"last_live"`
	AvatarURL    *string   `json:"avatarUrl,omitempty" db:"avatar_url"`
}
