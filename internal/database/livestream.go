package database

import (
	"strings"

	"github.com/galchammat/kadeem/internal/models"
)

func (db *DB) ListStreamers() ([]models.Streamer, error) {
	var streamers []models.Streamer
	rows, err := db.SQL.Query("SELECT id, name FROM streamers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s models.Streamer
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			return nil, err
		}
		streamers = append(streamers, s)
	}
	return streamers, nil
}

func (db *DB) ListStreams(filter *models.Stream) ([]models.Stream, error) {
	query := `SELECT * FROM streams`
	var where []string
	var args []interface{}

	if filter != nil && filter.ID != 0 {
		where = append(where, "id = ?")
		args = append(args, filter.ID)
	}
	if filter != nil && filter.StreamerID != 0 {
		where = append(where, "streamer_id = ?")
		args = append(args, filter.StreamerID)
	}
	if filter != nil && filter.Platform != "" {
		where = append(where, "platform = ?")
		args = append(args, filter.Platform)
	}
	if filter != nil && filter.ChannelName != "" {
		where = append(where, "channel_name = ?")
		args = append(args, filter.ChannelName)
	}
	if filter != nil && filter.ChannelID != "" {
		where = append(where, "channel_id = ?")
		args = append(args, filter.ChannelID)
	}

	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	rows, err := db.SQL.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var streams []models.Stream
	for rows.Next() {
		stream := models.Stream{}
		if err := rows.Scan(&stream.ID, &stream.StreamerID, &stream.Platform, &stream.ChannelName, &stream.ChannelID); err != nil {
			return nil, err
		}
		streams = append(streams, stream)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return streams, nil
}

func (db *DB) SaveStream(stream models.Stream) error {
	_, err := db.SQL.Exec(
		`INSERT INTO streams (streamer_id, platform, channel_name, channel_id) VALUES (?, ?, ?, ?)`,
		stream.StreamerID, stream.Platform, stream.ChannelName, stream.ChannelID,
	)
	return err
}
