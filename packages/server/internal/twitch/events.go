package twitch

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/model"
)

type hypeTrainItem struct {
	ID             string `json:"id"`
	EventType      string `json:"event_type"`
	EventTimestamp string `json:"event_timestamp"`
	EventData      struct {
		Level            int `json:"level"`
		Total            int `json:"total"`
		TopContributions []struct {
			UserName string `json:"user_name"`
			Type     string `json:"type"`
			Total    int    `json:"total"`
		} `json:"top_contributions"`
	} `json:"event_data"`
}

type clipItem struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	ViewCount   int     `json:"view_count"`
	CreatedAt   string  `json:"created_at"`
	CreatorName string  `json:"creator_name"`
	Duration    float64 `json:"duration"`
}

// FetchHypeTrainEvents fetches ended hype train events for the given broadcaster.
func (c *TwitchClient) FetchHypeTrainEvents(broadcasterID string) ([]model.StreamEvent, error) {
	params := url.Values{}
	params.Set("broadcaster_id", broadcasterID)
	params.Set("first", "100")

	response, statusCode, err := c.makeRequest("/helix/hypetrain/events?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("fetch hype train events: %w", err)
	}
	if statusCode != 200 {
		return nil, fmt.Errorf("fetch hype train events: unexpected status %d", statusCode)
	}

	var rawMessages []json.RawMessage
	if err := json.Unmarshal(response.Data, &rawMessages); err != nil {
		return nil, fmt.Errorf("unmarshal hype train events: %w", err)
	}

	var events []model.StreamEvent
	for _, raw := range rawMessages {
		var item hypeTrainItem
		if err := json.Unmarshal(raw, &item); err != nil {
			logging.Warn("failed to unmarshal hype train item", "error", err)
			continue
		}
		if item.EventType != "hypetrain.end" {
			continue
		}

		ts, err := time.Parse(time.RFC3339, item.EventTimestamp)
		if err != nil {
			logging.Warn("failed to parse hype train timestamp", "error", err)
			continue
		}

		level := strconv.Itoa(item.EventData.Level)
		description := fmt.Sprintf("Total %d pts · Level %d", item.EventData.Total, item.EventData.Level)
		if len(item.EventData.TopContributions) > 0 {
			top := item.EventData.TopContributions[0]
			description += fmt.Sprintf(" · Top: %s (%s)", top.UserName, top.Type)
		}

		externalID := item.ID
		events = append(events, model.StreamEvent{
			ChannelID:   broadcasterID,
			EventType:   model.StreamEventHypeTrain,
			Title:       fmt.Sprintf("Hype Train Level %d", item.EventData.Level),
			Description: description,
			Timestamp:   ts.Unix(),
			Value:       &level,
			ExternalID:  &externalID,
		})
	}
	return events, nil
}

// FetchTopClips fetches recent top clips for the given broadcaster.
func (c *TwitchClient) FetchTopClips(broadcasterID string) ([]model.StreamEvent, error) {
	params := url.Values{}
	params.Set("broadcaster_id", broadcasterID)
	params.Set("first", "20")

	response, statusCode, err := c.makeRequest("/helix/clips?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("fetch clips: %w", err)
	}
	if statusCode != 200 {
		return nil, fmt.Errorf("fetch clips: unexpected status %d", statusCode)
	}

	var rawMessages []json.RawMessage
	if err := json.Unmarshal(response.Data, &rawMessages); err != nil {
		return nil, fmt.Errorf("unmarshal clips: %w", err)
	}

	var events []model.StreamEvent
	for _, raw := range rawMessages {
		var item clipItem
		if err := json.Unmarshal(raw, &item); err != nil {
			logging.Warn("failed to unmarshal clip item", "error", err)
			continue
		}

		ts, err := time.Parse(time.RFC3339, item.CreatedAt)
		if err != nil {
			logging.Warn("failed to parse clip timestamp", "error", err)
			continue
		}

		views := strconv.Itoa(item.ViewCount)
		externalID := item.ID
		events = append(events, model.StreamEvent{
			ChannelID:   broadcasterID,
			EventType:   model.StreamEventClip,
			Title:       item.Title,
			Description: fmt.Sprintf("by %s · %d views · %.0fs", item.CreatorName, item.ViewCount, item.Duration),
			Timestamp:   ts.Unix(),
			Value:       &views,
			ExternalID:  &externalID,
		})
	}
	return events, nil
}
