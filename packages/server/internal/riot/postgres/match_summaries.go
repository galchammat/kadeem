package postgres

import (
	"context"
	"fmt"

	"github.com/galchammat/kadeem/internal/riot/models"
	"github.com/lib/pq"
)

func (s *DB) SaveMatchSummaryBatch(ctx context.Context, matchSummaries []models.MatchSummary) error {
	if len(matchSummaries) == 0 {
		return nil
	}

	ids := make([]int64, len(matchSummaries))
	regions := make([]string, len(matchSummaries))
	startedAts := make([]int64, len(matchSummaries))
	durations := make([]int, len(matchSummaries))
	queueIDs := make([]int, len(matchSummaries))

	for i, summary := range matchSummaries {
		ids[i] = summary.ID
		regions[i] = summary.Region
		startedAts[i] = summary.StartedAt
		durations[i] = summary.Duration
		queueIDs[i] = summary.QueueID
	}

	_, err := s.db.SQL.ExecContext(ctx, `
		INSERT INTO lol_matches (id, region, started_at, duration, queue_id)
		SELECT *
		FROM unnest(
			$1::bigint[],
			$2::text[],
			$3::bigint[],
			$4::integer[],
			$5::integer[]
		) AS summaries(id, region, started_at, duration, queue_id)
		ON CONFLICT (id, region) DO UPDATE SET
			started_at = EXCLUDED.started_at,
			duration = EXCLUDED.duration,
			queue_id = EXCLUDED.queue_id
	`, pq.Array(ids), pq.Array(regions), pq.Array(startedAts), pq.Array(durations), pq.Array(queueIDs))
	if err != nil {
		return fmt.Errorf("save match summary batch: %w", err)
	}

	return nil
}
