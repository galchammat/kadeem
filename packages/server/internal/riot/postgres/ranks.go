package postgres

import (
	"fmt"

	"github.com/galchammat/kadeem/internal/logging"
	riot "github.com/galchammat/kadeem/internal/riot/models"
)

func (s *DB) InsertPlayerRank(rank *riot.PlayerRank) error {
	query := `
        INSERT INTO player_ranks 
        (puuid, timestamp, tier, rank, league_points, wins, losses, queue_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        ON CONFLICT (puuid, timestamp, queue_id) DO UPDATE SET
            tier = EXCLUDED.tier,
            rank = EXCLUDED.rank,
            league_points = EXCLUDED.league_points,
            wins = EXCLUDED.wins,
            losses = EXCLUDED.losses`

	_, err := s.db.SQL.Exec(query,
		rank.PUUID, rank.Timestamp, rank.Tier, rank.Rank,
		rank.LeaguePoints, rank.Wins, rank.Losses, rank.QueueID)

	if err != nil {
		logging.Error("Failed to insert player rank", "puuid", rank.PUUID, "error", err)
	}
	return err
}

// GetRankAtTime fetches the rank closest to (but not after) a given timestamp
func (s *DB) GetRankAtTime(puuid string, queueID int, timestamp int64) (*riot.PlayerRank, error) {
	query := `
        SELECT puuid, timestamp, tier, rank, league_points, wins, losses, queue_id
        FROM player_ranks
        WHERE puuid = $1 AND queue_id = $2 AND timestamp <= $3
        ORDER BY timestamp DESC
        LIMIT 1`

	var rank riot.PlayerRank
	err := s.db.SQL.QueryRow(query, puuid, queueID, timestamp).Scan(
		&rank.PUUID, &rank.Timestamp, &rank.Tier, &rank.Rank,
		&rank.LeaguePoints, &rank.Wins, &rank.Losses, &rank.QueueID)

	if err != nil {
		logging.Error("Failed to get rank at time", "puuid", puuid, "queueID", queueID, "timestamp", timestamp, "error", err)
		return nil, fmt.Errorf("rank not found")
	}

	return &rank, nil
}
