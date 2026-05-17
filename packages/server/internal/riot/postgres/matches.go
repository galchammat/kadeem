package postgres

import (
	"context"
	"fmt"
)

func (s *DB) SaveMatchIDs(ctx context.Context, matchIDs []int64, region string) error {
	tx, err := s.db.SQL.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin save match ids: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck // rollback after commit is a no-op

	for _, matchID := range matchIDs {
		if err := ctx.Err(); err != nil {
			return err
		}

		_, err := tx.ExecContext(ctx, `
			INSERT INTO lol_matches (id, region, status, replay_status)
			VALUES ($1, $2, 'pending', 'pending')
			ON CONFLICT (id, region) DO NOTHING
		`, matchID, region)
		if err != nil {
			return fmt.Errorf("save match id %d: %w", matchID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit save match ids: %w", err)
	}

	return nil
}

func (s *DB) SaveMatchDetails(ctx context.Context, matchID int64, region string, startedAt int64, duration int, queueID int) error {
	_, err := s.db.SQL.ExecContext(ctx, `
		UPDATE lol_matches
		SET started_at = $1,
			duration = $2,
			queue_id = $3
		WHERE id = $4
		  AND region = $5
	`, startedAt, duration, queueID, matchID, region)
	return err
}
