package model

// StreamEventType classifies a stream event.
type StreamEventType string

const (
	StreamEventHypeTrain StreamEventType = "hype_train"
	StreamEventClip      StreamEventType = "clip"
)

// StreamEvent is a notable moment that occurred during a live broadcast.
type StreamEvent struct {
	ID          int64           `json:"id"`
	ChannelID   string          `json:"channel_id"`
	EventType   StreamEventType `json:"type"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Timestamp   int64           `json:"timestamp"` // Unix seconds
	Value       *string         `json:"value,omitempty"`
	ExternalID  *string         `json:"-"` // Twitch event ID; used for deduplication
}

// StreamEventFilter constrains a ListStreamEvents query.
type StreamEventFilter struct {
	ChannelID    *string
	StreamerID   *int64
	EventType    *StreamEventType
	TimestampMin *int64
	TimestampMax *int64
}
