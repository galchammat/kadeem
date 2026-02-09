package database

import (
	"fmt"
	"strings"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/model"
)

func (db *DB) ListBroadcasts(filter *model.Broadcast, limit int, offset int) ([]model.Broadcast, error) {
	if filter == nil || filter.ChannelID == "" {
		return nil, fmt.Errorf("channel_id is required for ListBroadcasts")
	}

	query := `SELECT id, channel_id, title, url, thumbnail_url, viewable, created_at, published_at, duration
		FROM broadcasts WHERE channel_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := db.SQL.Query(query, filter.ChannelID, limit, offset)
	if err != nil {
		logging.Error("Failed to query broadcasts", "channelID", filter.ChannelID, "error", err)
		return nil, err
	}
	defer rows.Close()

	var broadcasts []model.Broadcast
	for rows.Next() {
		var b model.Broadcast
		if err := rows.Scan(
			&b.ID,
			&b.ChannelID,
			&b.Title,
			&b.URL,
			&b.ThumbnailURL,
			&b.Viewable,
			&b.CreatedAt,
			&b.PublishedAt,
			&b.Duration,
		); err != nil {
			logging.Error("Failed to scan broadcast row", "error", err)
			return nil, err
		}
		broadcasts = append(broadcasts, b)
	}
	if err := rows.Err(); err != nil {
		logging.Error("Error iterating over broadcast rows", "error", err)
		return nil, err
	}
	return broadcasts, nil
}

func (db *DB) InsertBroadcasts(broadcasts []model.Broadcast) error {
	if len(broadcasts) == 0 {
		return nil
	}

	query := `INSERT INTO broadcasts (
			channel_id, title, url, thumbnail_url, viewable, created_at, published_at, duration
		) VALUES `

	var (
		placeholders []string
		args         []any
	)

	argN := 1
	for i, b := range broadcasts {
		if b.ChannelID == "" {
			logging.Warn("InsertBroadcasts: missing required field in broadcast", "index", i, "broadcast", b)
			continue
		}
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			argN, argN+1, argN+2, argN+3, argN+4, argN+5, argN+6, argN+7))
		args = append(args,
			b.ChannelID,
			b.Title,
			b.URL,
			b.ThumbnailURL,
			b.Viewable,
			b.CreatedAt,
			b.PublishedAt,
			b.Duration,
		)
		argN += 8
	}

	if len(placeholders) == 0 {
		return nil
	}

	query += strings.Join(placeholders, ", ")
	query += ` ON CONFLICT (channel_id, url) DO NOTHING`

	_, err := db.SQL.Exec(query, args...)
	if err != nil {
		logging.Error("Failed to insert broadcasts", "error", err)
	}
	return err
}
