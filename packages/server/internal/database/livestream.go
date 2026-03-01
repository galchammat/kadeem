package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/model"
)

// SaveStreamer saves a streamer to the database (shared pool)
func (db *DB) SaveStreamer(streamer model.Streamer) (int64, error) {
	var id int64
	err := db.SQL.QueryRow(
		`INSERT INTO streamers (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id`,
		streamer.Name,
	).Scan(&id)
	if err != nil {
		logging.Error("Failed to save streamer to database", "name", streamer.Name, "error", err)
		return 0, err
	}
	return id, nil
}

// GetStreamerByName retrieves a streamer by name
func (db *DB) GetStreamerByName(name string) (*model.Streamer, error) {
	var streamer model.Streamer
	err := db.SQL.QueryRow(`SELECT id, name FROM streamers WHERE name = $1`, name).Scan(&streamer.ID, &streamer.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		logging.Error("Failed to get streamer", "name", name, "error", err)
		return nil, err
	}
	return &streamer, nil
}

// GetStreamerByID retrieves a streamer by ID
func (db *DB) GetStreamerByID(id int) (*model.Streamer, error) {
	var streamer model.Streamer
	err := db.SQL.QueryRow(`SELECT id, name FROM streamers WHERE id = $1`, id).Scan(&streamer.ID, &streamer.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		logging.Error("Failed to get streamer", "id", id, "error", err)
		return nil, err
	}
	return &streamer, nil
}

// FindOrCreateStreamer finds or creates a streamer (idempotent)
func (db *DB) FindOrCreateStreamer(name string) (*model.Streamer, error) {
	streamer, err := db.GetStreamerByName(name)
	if err != nil {
		return nil, err
	}
	if streamer != nil {
		return streamer, nil
	}

	id, err := db.SaveStreamer(model.Streamer{Name: name})
	if err != nil {
		return nil, err
	}

	return &model.Streamer{ID: id, Name: name}, nil
}

// ListTrackedStreamers returns streamers a user is tracking with pagination
func (db *DB) ListTrackedStreamers(userID string, limit, offset int) ([]model.Streamer, error) {
	query := `SELECT s.id, s.name 
	          FROM streamers s
	          INNER JOIN user_tracked_streamers uts ON s.id = uts.streamer_id
	          WHERE uts.user_id = $1
	          ORDER BY uts.tracked_at DESC
	          LIMIT $2 OFFSET $3`

	rows, err := db.SQL.Query(query, userID, limit, offset)
	if err != nil {
		logging.Error("Failed to list tracked streamers", "userID", userID, "error", err)
		return nil, err
	}
	defer rows.Close()

	var streamers []model.Streamer
	for rows.Next() {
		var s model.Streamer
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			logging.Error("Failed to scan tracked streamer row", "error", err)
			return nil, err
		}
		streamers = append(streamers, s)
	}
	if err := rows.Err(); err != nil {
		logging.Error("Error iterating over tracked streamer rows", "error", err)
		return nil, err
	}
	return streamers, nil
}

// TrackStreamer adds a tracking relationship (idempotent)
func (db *DB) TrackStreamer(userID string, streamerID int64) error {
	query := `INSERT INTO user_tracked_streamers (user_id, streamer_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := db.SQL.Exec(query, userID, streamerID)
	if err != nil {
		logging.Error("Failed to track streamer", "userID", userID, "streamerID", streamerID, "error", err)
	}
	return err
}

// UntrackStreamer removes a tracking relationship
func (db *DB) UntrackStreamer(userID string, streamerID int64) error {
	query := `DELETE FROM user_tracked_streamers WHERE user_id = $1 AND streamer_id = $2`
	_, err := db.SQL.Exec(query, userID, streamerID)
	if err != nil {
		logging.Error("Failed to untrack streamer", "userID", userID, "streamerID", streamerID, "error", err)
	}
	return err
}

// IsTrackingStreamer checks if user is tracking a streamer
func (db *DB) IsTrackingStreamer(userID string, streamerID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM user_tracked_streamers WHERE user_id = $1 AND streamer_id = $2)`
	var exists bool
	err := db.SQL.QueryRow(query, userID, streamerID).Scan(&exists)
	if err != nil {
		logging.Error("Failed to check streamer tracking status", "userID", userID, "streamerID", streamerID, "error", err)
		return false, err
	}
	return exists, nil
}

// GetTrackedStreamersForSync returns all streamers with at least one tracker (for background jobs)
func (db *DB) GetTrackedStreamersForSync() ([]model.Streamer, error) {
	query := `SELECT DISTINCT s.id, s.name 
	          FROM streamers s
	          INNER JOIN user_tracked_streamers uts ON s.id = uts.streamer_id`

	rows, err := db.SQL.Query(query)
	if err != nil {
		logging.Error("Failed to get tracked streamers for sync", "error", err)
		return nil, err
	}
	defer rows.Close()

	var streamers []model.Streamer
	for rows.Next() {
		var s model.Streamer
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			logging.Error("Failed to scan streamer row for sync", "error", err)
			return nil, err
		}
		streamers = append(streamers, s)
	}
	if err := rows.Err(); err != nil {
		logging.Error("Error iterating over streamers for sync", "error", err)
		return nil, err
	}
	return streamers, nil
}

// DeleteStreamer deletes a streamer by name (admin only)
func (db *DB) DeleteStreamer(name string) (bool, error) {
	res, err := db.SQL.Exec(
		`DELETE FROM streamers WHERE name = $1`,
		name,
	)
	if err != nil {
		logging.Error("Failed to delete streamer from database", "name", name, "error", err)
		return false, err
	}
	n, _ := res.RowsAffected()
	return (n != 0), nil
}

// ListStreamers lists all streamers with pagination (for admin/internal use)
func (db *DB) ListStreamers(limit, offset int) ([]model.Streamer, error) {
	rows, err := db.SQL.Query("SELECT id, name FROM streamers ORDER BY name LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		logging.Error("Failed to query streamers from database", "error", err)
		return nil, err
	}
	defer rows.Close()

	var streamers []model.Streamer
	for rows.Next() {
		var s model.Streamer
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			logging.Error("Failed to scan streamer row", "error", err)
			return nil, err
		}
		streamers = append(streamers, s)
	}
	if err := rows.Err(); err != nil {
		logging.Error("Error iterating over streamer rows", "error", err)
		return nil, err
	}
	return streamers, nil
}

// ListChannels lists channels with optional filtering and pagination
func (db *DB) ListChannels(filter *model.ChannelFilter, limit, offset int) ([]model.Channel, error) {
	query := `SELECT id, streamer_id, platform, channel_name, avatar_url, synced_at FROM channels`
	var where []string
	var args []any
	argN := 1

	if filter != nil {
		if filter.ID != nil {
			where = append(where, fmt.Sprintf("id = $%d", argN))
			args = append(args, *filter.ID)
			argN++
		}
		if filter.StreamerID != nil {
			where = append(where, fmt.Sprintf("streamer_id = $%d", argN))
			args = append(args, *filter.StreamerID)
			argN++
		}
		if filter.Platform != nil {
			where = append(where, fmt.Sprintf("platform = $%d", argN))
			args = append(args, *filter.Platform)
			argN++
		}
		if filter.ChannelName != nil {
			where = append(where, fmt.Sprintf("channel_name LIKE $%d", argN))
			args = append(args, *filter.ChannelName)
			argN++
		}
	}

	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}
	query += fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", argN, argN+1)
	args = append(args, limit, offset)

	rows, err := db.SQL.Query(query, args...)
	if err != nil {
		logging.Error("Failed to query channels from database", "error", err)
		return nil, err
	}
	defer rows.Close()

	var channels []model.Channel
	for rows.Next() {
		var ch model.Channel
		var syncedAt sql.NullTime
		if err := rows.Scan(&ch.ID, &ch.StreamerID, &ch.Platform, &ch.ChannelName, &ch.AvatarURL, &syncedAt); err != nil {
			logging.Error("Failed to scan channel row", "error", err)
			return nil, err
		}
		if syncedAt.Valid {
			unixTime := syncedAt.Time.Unix()
			ch.SyncedAt = &unixTime
		}
		channels = append(channels, ch)
	}
	if err := rows.Err(); err != nil {
		logging.Error("Error iterating over channel rows", "error", err)
		return nil, err
	}
	return channels, nil
}

func (db *DB) SaveChannel(channel model.Channel) (bool, error) {
	res, err := db.SQL.Exec(
		`INSERT INTO channels (streamer_id, platform, channel_name, id, avatar_url) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO NOTHING`,
		channel.StreamerID, channel.Platform, channel.ChannelName, channel.ID, channel.AvatarURL,
	)
	if err != nil {
		logging.Error("Failed to save channel to database", "channelID", channel.ID, "channelName", channel.ChannelName, "error", err)
		return false, err
	}
	n, _ := res.RowsAffected()
	return (n != 0), nil
}

// allowedChannelColumns is the set of columns that can be updated via UpdateChannel.
var allowedChannelColumns = map[string]bool{
	"channel_name": true,
	"avatar_url":   true,
	"synced_at":    true,
	"platform":     true,
}

func (db *DB) UpdateChannel(channelID string, updates map[string]any) (bool, error) {
	var setClauses []string
	var args []any
	argN := 1

	for column, value := range updates {
		if !allowedChannelColumns[column] {
			return false, fmt.Errorf("disallowed column: %s", column)
		}
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, argN))
		args = append(args, value)
		argN++
	}
	if len(setClauses) == 0 {
		return false, nil
	}
	args = append(args, channelID)

	query := `UPDATE channels SET ` + strings.Join(setClauses, ", ") + fmt.Sprintf(` WHERE id = $%d`, argN)

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
		`DELETE FROM channels WHERE id = $1`,
		channelID,
	)
	if err != nil {
		logging.Error("Failed to delete channel from database", "channelID", channelID, "error", err)
		return false, err
	}
	n, _ := res.RowsAffected()
	return (n != 0), nil
}
