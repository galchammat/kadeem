package models

import (
	"encoding/json"
	"time"
)

type ChannelSearchResponse []struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"thumbnail_url"`
}

type TwitchResponse struct {
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

// type BroadcastListResponse struct {
// 	Data []struct {
// 		ID           string `json:"id"`
// 		Title        string `json:"title"`
// 		CreatedAt    string `json:"created_at"`
// 		PublishedAt  string `json:"published_at"`
// 		URL          string `json:"url"`
// 		ThumbnailURL string `json:"thumbnail_url"`
// 		Viewable     string `json:"viewable"`
// 		Duration     string `json:"duration"`
// 	} `json:"data"`
// }
