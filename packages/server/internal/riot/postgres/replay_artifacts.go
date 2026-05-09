package postgres

import (
	"context"
	"fmt"
	"strconv"

	"github.com/galchammat/kadeem/internal/syncer"
)

var _ syncer.ArtifactOps = replayArtifacts{}

type replayArtifacts struct {
	db *DB
}

func (db *DB) ReplayArtifactOps() syncer.ArtifactOps {
	return replayArtifacts{db: db}
}

func (s replayArtifacts) ClaimPending(ctx context.Context, limit int) ([]syncer.Artifact, error) {
	if limit <= 0 {
		return nil, nil
	}

	rows, err := s.db.db.SQL.QueryContext(ctx, `
		WITH pending AS (
			SELECT id
			FROM lol_matches
			WHERE replay_s3_key IS NULL
			ORDER BY id DESC
			LIMIT $1
			FOR UPDATE SKIP LOCKED
		)
		UPDATE lol_matches
		SET replay_sync_attempted_at = NOW(),
			replay_sync_error = NULL
		FROM pending
		WHERE lol_matches.id = pending.id
		RETURNING lol_matches.id
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("claim pending replays: %w", err)
	}
	defer rows.Close()

	artifacts := make([]syncer.Artifact, 0, limit)
	for rows.Next() {
		var matchID int64
		if err := rows.Scan(&matchID); err != nil {
			return nil, fmt.Errorf("scan replay artifact: %w", err)
		}

		id := strconv.FormatInt(matchID, 10)
		artifacts = append(artifacts, syncer.Artifact{
			ID:         id,
			ExternalID: id,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate replay artifacts: %w", err)
	}

	return artifacts, nil
}

func (s replayArtifacts) MarkDone(ctx context.Context, id string, s3Key string) error {
	matchID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("parse replay artifact id %q: %w", id, err)
	}

	result, err := s.db.db.SQL.ExecContext(ctx, `
		UPDATE lol_matches
		SET replay_s3_key = $2,
			replay_sync_error = NULL
		WHERE id = $1
	`, matchID, s3Key)
	if err != nil {
		return fmt.Errorf("mark replay done %q: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read replay mark done result %q: %w", id, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("replay artifact %q not found", id)
	}

	return nil
}

func (s replayArtifacts) MarkFailed(ctx context.Context, id string, err error) error {
	matchID, parseErr := strconv.ParseInt(id, 10, 64)
	if parseErr != nil {
		return fmt.Errorf("parse replay artifact id %q: %w", id, parseErr)
	}

	_, updateErr := s.db.db.SQL.ExecContext(ctx, `
		UPDATE lol_matches
		SET replay_sync_error = $2,
			replay_sync_attempted_at = NOW()
		WHERE id = $1
	`, matchID, err.Error())
	if updateErr != nil {
		return fmt.Errorf("mark replay failed %q: %w", id, updateErr)
	}

	return nil
}
