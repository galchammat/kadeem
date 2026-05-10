package postgres

import (
	"context"
	"fmt"

	"github.com/galchammat/kadeem/internal/riot/models"
)

func (s *DB) ClaimPendingMatch(ctx context.Context) (*int64, *string, error) {
	var matchID *int64
	var region *string

	err := s.db.SQL.QueryRowContext(ctx, `
	WITH claimed AS (
		SELECT id
		FROM matches
		WHERE status = 'pending'
		FOR UPDATE SKIP LOCKED
		LIMIT 1
	)
	UPDATE matches m
	SET status = 'processing'
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
		SET status = 'completed'
		WHERE id = $1
		  AND region = $2
	`, matchId, region)

	return err
}

func (s *DB) NackMatch(ctx context.Context, matchId int64, region string) error {
	_, err := s.db.SQL.ExecContext(ctx, `
		UPDATE lol_matches
		SET status = 'failed'
		WHERE id = $1
		  AND region = $2
	`, matchId, region)

	return err
}

func (s *DB) SaveMatchDetails(ctx context.Context, match models.Match) error {
	return nil
}
