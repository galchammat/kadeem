package database

import (
	"strings"

	"github.com/galchammat/kadeem/internal/models"
)

func (db *DB) SaveStreamer(streamer models.Streamer) (bool, error) {
	res, err := db.SQL.Exec(
		`INSERT OR IGNORE INTO streamers (name) VALUES (?)`,
		streamer.Name,
	)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return (n != 0), nil
}

func (db *DB) DeleteStreamer(name string) (bool, error) {
	res, err := db.SQL.Exec(
		`DELETE FROM streamers WHERE name = ?`,
		name,
	)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return (n != 0), nil
}

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

func (db *DB) ListChannels(filter *models.Channel) ([]models.Channel, error) {
	query := `SELECT * FROM channels`
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

	var streams []models.Channel
	for rows.Next() {
		stream := models.Channel{}
		if err := rows.Scan(&stream.ID, &stream.StreamerID, &stream.Platform, &stream.ChannelName, &stream.ChannelID, &stream.AvatarURL); err != nil {
			return nil, err
		}
		streams = append(streams, stream)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return streams, nil
}

func (db *DB) SaveChannel(channel models.Channel) (bool, error) {
	res, err := db.SQL.Exec(
		`INSERT OR IGNORE INTO channels (streamer_id, platform, channel_name, channel_id, avatar_url) VALUES (?, ?, ?, ?, ?)`,
		channel.StreamerID, channel.Platform, channel.ChannelName, channel.ChannelID, channel.AvatarURL,
	)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return (n != 0), nil
}
