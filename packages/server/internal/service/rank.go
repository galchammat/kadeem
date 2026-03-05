package service

import (
	"fmt"
	"time"

	"github.com/galchammat/kadeem/internal/database"
	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/model"
	riot "github.com/galchammat/kadeem/internal/riot/api"
)

type RankService struct {
	db   *database.DB
	riot *riot.Client
}

func NewRankService(db *database.DB, riot *riot.Client) *RankService {
	return &RankService{db: db, riot: riot}
}

// SyncRank fetches current rank for an account and stores a snapshot.
func (s *RankService) SyncRank(account *model.LolAccount) error {
	summonerID, err := s.riot.FetchSummonerID(account.PUUID, account.Region)
	if err != nil {
		return err
	}

	entries, err := s.riot.FetchRankEntries(summonerID, account.Region)
	if err != nil {
		return err
	}

	timestamp := time.Now().Unix()
	for _, entry := range entries {
		queueID := queueTypeToID(entry.QueueType)
		if queueID == 0 {
			continue
		}

		rank := &model.PlayerRank{
			PUUID:        account.PUUID,
			Timestamp:    timestamp,
			Tier:         entry.Tier,
			Rank:         entry.Rank,
			LeaguePoints: entry.LeaguePoints,
			Wins:         entry.Wins,
			Losses:       entry.Losses,
			QueueId:      queueID,
		}

		if err := s.db.InsertPlayerRank(rank); err != nil {
			logging.Error("Failed to insert rank", "puuid", account.PUUID, "queueID", queueID, "error", err)
			return fmt.Errorf("failed to insert rank: %w", err)
		}
	}

	logging.Debug("Synced rank for account", "puuid", account.PUUID, "numEntries", len(entries))
	return nil
}

func queueTypeToID(queueType string) int {
	switch queueType {
	case "RANKED_SOLO_5x5":
		return 420
	case "RANKED_FLEX_SR":
		return 440
	default:
		return 0
	}
}
