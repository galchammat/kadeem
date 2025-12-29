package database

import (
	"fmt"
	"strings"

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
	if filter.URL != "" {
		where = append(where, "url = ?")
		args = append(args, filter.URL)
	}
	if filter.StartedAt != 0 {
		where = append(where, "started_at = ?")
		args = append(args, filter.StartedAt)
	}
	if filter.EndedAt != 0 {
		where = append(where, "ended_at = ?")
		args = append(args, filter.EndedAt)
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
		if err := rows.Scan(&b.ID, &b.ChannelID, &b.URL, &b.StartedAt, &b.EndedAt); err != nil {
			return nil, err
		}
		broadcasts = append(broadcasts, b)
	}

	return broadcasts, nil
}
