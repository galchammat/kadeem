package database

import (
	"fmt"
	"strings"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/models"
)

func (db *DB) ListBroadcasts(filter *models.BroadcastFilter, limit int, offset *int) ([]models.Broadcast, error) {
	query := `SELECT * FROM broadcasts`

	// ChannelID is required
	if filter == nil || filter.ChannelID == nil {
		return []models.Broadcast{}, fmt.Errorf("ChannelID must be specified when calling ListBroadcasts.")
	}

	// Build WHERE clauses using BuildQueryArgs
	whereClauses, args, err := db.BuildQueryArgs(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to build query args: %w", err)
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

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

	query += strings.Join(placeholders, ", ")

	_, err := db.SQL.Exec(query, args...)

	return err
}
