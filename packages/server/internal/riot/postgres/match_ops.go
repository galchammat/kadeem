package postgres

import (
	"context"
	"fmt"
)

func (s *DB) ClaimPendingMatch(ctx context.Context) (*int64, *string, error) {
	var matchID *int64
	var region *string

	err := s.db.SQL.QueryRowContext(ctx, `
	WITH claimed AS (
		SELECT id
		FROM lol_matches
		WHERE status = 'pending'
		FOR UPDATE SKIP LOCKED
		LIMIT 1
	)
	UPDATE lol_matches m
	SET status = 'processing',
		updated_at = NOW()
	FROM claimed
	WHERE m.id = claimed.id
	RETURNING m.id, m.region
`).Scan(&matchID, &region)
	if err != nil {
		return nil, nil, fmt.Errorf("ClaimPendingMatch failed: %w", err)
	}

	return matchID, region, nil
}

func (s *DB) AckMatch(ctx context.Context, matchId int64, region string) error {
	_, err := s.db.SQL.ExecContext(ctx, `
		UPDATE lol_matches
		SET status = 'completed',
			updated_at = NOW()
		WHERE id = $1
		  AND region = $2
	`, matchId, region)

	return err
}

func (s *DB) NackMatch(ctx context.Context, matchId int64, region string) error {
	_, err := s.db.SQL.ExecContext(ctx, `
		UPDATE lol_matches
		SET status = 'failed',
			updated_at = NOW()
		WHERE id = $1
		  AND region = $2
	`, matchId, region)

	return err
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
