package database

import (
	"database/sql"
	"strings"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/models"
)

func (db *DB) SaveStreamer(streamer models.Streamer) (bool, error) {
	res, err := db.SQL.Exec(
		`INSERT OR IGNORE INTO streamers (name) VALUES (?)`,
		streamer.Name,
	)
	if err != nil {
		logging.Error("Failed to save streamer to database", "name", streamer.Name, "error", err)
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
		logging.Error("Failed to delete streamer from database", "name", name, "error", err)
		return false, err
	}
	n, _ := res.RowsAffected()
	return (n != 0), nil
}

func (db *DB) ListStreamers() ([]models.Streamer, error) {
	var streamers []models.Streamer
	rows, err := db.SQL.Query("SELECT id, name FROM streamers")
	if err != nil {
		logging.Error("Failed to query streamers from database", "error", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s models.Streamer
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			logging.Error("Failed to scan streamer row", "error", err)
			return nil, err
		}
		streamers = append(streamers, s)
	}
	return streamers, nil
}

func (db *DB) ListChannels(filter *models.ChannelFilter) ([]models.Channel, error) {
	query := `SELECT * FROM channels`

	// Build WHERE clauses using BuildQueryArgs
	whereClauses, args, err := db.BuildQueryArgs(filter)
	if err != nil {
		logging.Error("Failed to build query args for ListChannels", "error", err)
		return nil, err
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	rows, err := db.SQL.Query(query, args...)
	if err != nil {
		logging.Error("Failed to query channels from database", "error", err)
		return nil, err
	}
	defer rows.Close()

	var channels []models.Channel
	for rows.Next() {
		channel := models.Channel{}
		var syncedAt sql.NullTime
		if err := rows.Scan(&channel.ID, &channel.StreamerID, &channel.Platform, &channel.ChannelName, &channel.AvatarURL, &syncedAt); err != nil {
			logging.Error("Failed to scan channel row", "error", err)
			return nil, err
		}

		if syncedAt.Valid {
			unixTime := syncedAt.Time.Unix()
			channel.SyncedAt = &unixTime
		}
		channels = append(channels, channel)
	}
	if err := rows.Err(); err != nil {
		logging.Error("Error iterating over channel rows", "error", err)
		return nil, err
	}
	return channels, nil
}

func (db *DB) SaveChannel(channel models.Channel) (bool, error) {
	res, err := db.SQL.Exec(
		`INSERT OR IGNORE INTO channels (streamer_id, platform, channel_name, id, avatar_url) VALUES (?, ?, ?, ?, ?)`,
		channel.StreamerID, channel.Platform, channel.ChannelName, channel.ID, channel.AvatarURL,
	)
	if err != nil {
		logging.Error("Failed to save channel to database", "channelID", channel.ID, "channelName", channel.ChannelName, "error", err)
		return false, err
	}
	n, _ := res.RowsAffected()
	return (n != 0), nil
}

func (db *DB) UpdateChannel(channelID string, updates map[string]interface{}) (bool, error) {
	var setClauses []string
	var args []interface{}

	for column, value := range updates {
		setClauses = append(setClauses, column+" = ?")
		args = append(args, value)
	}
	args = append(args, channelID)

	query := `UPDATE channels SET ` + strings.Join(setClauses, ", ") + ` WHERE id = ?`

	res, err := db.SQL.Exec(query, args...)
	if err != nil {
		logging.Error("Failed to update channel in database", "channelID", channelID, "error", err)
		return false, err
	}
	n, _ := res.RowsAffected()
	return (n != 0), nil
}

func (db *DB) DeleteChannel(channelID string) (bool, error) {
	res, err := db.SQL.Exec(
		`DELETE FROM channels WHERE id = ?`,
		channelID,
	)
	if err != nil {
		logging.Error("Failed to delete channel from database", "channelID", channelID, "error", err)
		return false, err
	}
	n, _ := res.RowsAffected()
	return (n != 0), nil
}
