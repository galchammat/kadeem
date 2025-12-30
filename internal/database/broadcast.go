package database

import (
	"fmt"
	"strings"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/models"
)

func (db *DB) ListBroadcasts(filter *models.Broadcast, limit int, offset *int) ([]models.Broadcast, error) {
	query := `SELECT * FROM broadcasts`
	var where []string
	var args []interface{}

	if filter != nil && filter.ChannelID != "" {
		where = append(where, "channel_id = ?")
		args = append(args, filter.ChannelID)
	} else {
		return []models.Broadcast{}, fmt.Errorf("ChannelID must be specified when calling ListBroadcasts.")
	}

	if filter.ID != 0 {
		where = append(where, "id = ?")
		args = append(args, filter.ID)
	}
	if filter.Title != "" {
		where = append(where, "title = ?")
		args = append(args, filter.Title)
	}
	if filter.URL != "" {
		where = append(where, "url = ?")
		args = append(args, filter.URL)
	}
	if filter.ThumbnailURL != "" {
		where = append(where, "thumbnail_url = ?")
		args = append(args, filter.ThumbnailURL)
	}
	if filter.Viewable != "" {
		where = append(where, "viewable = ?")
		args = append(args, filter.Viewable)
	}
	if filter.CreatedAt != 0 {
		where = append(where, "created_at = ?")
		args = append(args, filter.CreatedAt)
	}
	if filter.PublishedAt != 0 {
		where = append(where, "published_at = ?")
		args = append(args, filter.PublishedAt)
	}
	if filter.Duration != 0 {
		where = append(where, "duration = ?")
		args = append(args, filter.Duration)
	}

	query += " WHERE " + strings.Join(where, " AND ")

	rows, err := db.SQL.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var broadcasts []models.Broadcast
	for rows.Next() {
		var b models.Broadcast
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
			return nil, err
		}
		broadcasts = append(broadcasts, b)
	}

	return broadcasts, nil
}

func (db *DB) InsertBroadcasts(broadcasts []models.Broadcast) error {
	if len(broadcasts) == 0 {
		return nil
	}

	query := `INSERT INTO broadcasts (
			channel_id, title, url, thumbnail_url, viewable, created_at, published_at, duration
		) VALUES `

	var (
		placeholders []string
		args         []interface{}
	)

	for i, b := range broadcasts {
		if b.ChannelID == "" {
			logging.Warn("InsertBroadcasts: missing required field in broadcast", "index", i, "broadcast", b)
			continue
		}
		placeholders = append(placeholders, "(?, ?, ?, ?, ?, ?, ?, ?)")
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
	}

	logging.Debug("args", "args", args)

	query += strings.Join(placeholders, ", ")

	_, err := db.SQL.Exec(query, args...)

	return err
}
