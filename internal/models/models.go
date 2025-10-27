package models

type Streamer struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type StreamID int64

type Channel struct {
	ID          StreamID `json:"id" db:"id"`
	StreamerID  int64    `json:"streamerId" db:"streamer_id"`
	Platform    string   `json:"platform" db:"platform"`
	ChannelName string   `json:"channelName" db:"channel_name"`
	ChannelID   string   `json:"channelId" db:"channel_id"`
	AvatarURL   string   `json:"avatarUrl" db:"avatar_url"`
}

type Broadcast struct {
	ID        int64    `db:"id"`
	StreamID  StreamID `db:"stream_id"`
	URL       string   `db:"url"`
	StartedAt int64    `db:"started_at"`
	EndedAt   int64    `db:"ended_at,omitempty"`
}

type StreamerView struct {
	StreamerID   int64     `json:"id" db:"id"`
	StreamerName string    `json:"name" db:"name"`
	Channels     []Channel `json:"channels" db:"streams"`
	LastLive     *int64    `json:"lastLive,omitempty" db:"last_live"`
	AvatarURL    *string   `json:"avatarUrl,omitempty" db:"avatar_url"`
}
