package models

type Streamer struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type StreamID int64

type Channel struct {
	ID          StreamID `db:"id"`
	StreamerID  int64    `db:"streamer_id"`
	Platform    string   `db:"platform"`
	ChannelName string   `json:"display_name" db:"channel_name"`
	ChannelID   string   `json:"id" db:"channel_id"`
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
	Channels     []Channel `json:"streams" db:"streams"`
	LastLive     *int64    `json:"last_live,omitempty" db:"last_live"`
}
