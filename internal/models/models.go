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

type StreamerView struct {
	StreamerID   int64     `json:"id" db:"id"`
	StreamerName string    `json:"name" db:"name"`
	Channels     []Channel `json:"channels" db:"streams"`
	LastLive     *int64    `json:"lastLive,omitempty" db:"last_live"`
	AvatarURL    *string   `json:"avatarUrl,omitempty" db:"avatar_url"`
}
