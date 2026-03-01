package database

import (
	"fmt"
	"strings"

	"github.com/galchammat/kadeem/internal/model"
)

// ListStreamEvents returns stream events matching the filter, ordered by timestamp descending.
func (db *DB) ListStreamEvents(filter *model.StreamEventFilter, limit, offset int) ([]model.StreamEvent, error) {
	query := `SELECT se.id, se.channel_id, se.event_type, se.title, se.description,
	                 se.timestamp, se.value, se.external_id
	          FROM stream_events se`

	var joins []string
	var where []string
	var args []any
	argN := 1

	if filter != nil {
		if filter.StreamerID != nil {
			joins = append(joins, "INNER JOIN channels c ON se.channel_id = c.id")
			where = append(where, fmt.Sprintf("c.streamer_id = $%d", argN))
			args = append(args, *filter.StreamerID)
			argN++
		}
		if filter.ChannelID != nil {
			where = append(where, fmt.Sprintf("se.channel_id = $%d", argN))
			args = append(args, *filter.ChannelID)
			argN++
		}
		if filter.EventType != nil {
			where = append(where, fmt.Sprintf("se.event_type = $%d", argN))
			args = append(args, string(*filter.EventType))
			argN++
		}
		if filter.TimestampMin != nil {
			where = append(where, fmt.Sprintf("se.timestamp >= $%d", argN))
			args = append(args, *filter.TimestampMin)
			argN++
		}
		if filter.TimestampMax != nil {
			where = append(where, fmt.Sprintf("se.timestamp <= $%d", argN))
			args = append(args, *filter.TimestampMax)
			argN++
		}
	}

	if len(joins) > 0 {
		query += " " + strings.Join(joins, " ")
	}
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}
	query += fmt.Sprintf(" ORDER BY se.timestamp DESC LIMIT $%d OFFSET $%d", argN, argN+1)
	args = append(args, limit, offset)

	rows, err := db.SQL.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list stream events: %w", err)
	}
	defer rows.Close()

	var events []model.StreamEvent
	for rows.Next() {
		var e model.StreamEvent
		if err := rows.Scan(
			&e.ID, &e.ChannelID, &e.EventType, &e.Title, &e.Description,
			&e.Timestamp, &e.Value, &e.ExternalID,
		); err != nil {
			return nil, fmt.Errorf("scan stream event: %w", err)
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

// UpsertStreamEvents inserts stream events, ignoring duplicates by (channel_id, external_id).
func (db *DB) UpsertStreamEvents(events []model.StreamEvent) error {
	if len(events) == 0 {
		return nil
	}

	tx, err := db.SQL.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	stmt, err := tx.Prepare(`
		INSERT INTO stream_events (channel_id, event_type, title, description, timestamp, value, external_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (channel_id, external_id)
		WHERE external_id IS NOT NULL
		DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, e := range events {
		if _, err := stmt.Exec(
			e.ChannelID, string(e.EventType), e.Title, e.Description,
			e.Timestamp, e.Value, e.ExternalID,
		); err != nil {
			return fmt.Errorf("upsert stream event %v: %w", e.ExternalID, err)
		}
	}

	return tx.Commit()
}
