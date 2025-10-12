package models

type Streamer struct {
	Name    string   `json:"name" db:"name"`
	Streams []Stream `json:"streams,omitempty" db:"-"`
}

type StreamID int64

type Stream struct {
	ID       StreamID
	Platform string
	Username string
	Streamer string
}

type Broadcast struct {
	ID        int64
	StreamID  StreamID
	StartTime int64
	EndTime   int64
}
