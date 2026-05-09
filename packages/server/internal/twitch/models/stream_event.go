package models

type StreamEventType string

const (
	StreamEventHypeTrain StreamEventType = "hype_train"
	StreamEventClip      StreamEventType = "clip"
)

type StreamEvent struct {
	ID          int64           `json:"id"`
	ChannelID   string          `json:"channel_id"`
	EventType   StreamEventType `json:"type"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Timestamp   int64           `json:"timestamp"`
	Value       *string         `json:"value,omitempty"`
	ExternalID  *string         `json:"-"`
}

type StreamEventFilter struct {
	ChannelID    *string
	StreamerID   *int64
	EventType    *StreamEventType
	TimestampMin *int64
	TimestampMax *int64
}
