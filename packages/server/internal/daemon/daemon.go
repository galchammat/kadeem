package daemon

import (
	"context"
	"os"
	"time"

	"github.com/galchammat/kadeem/internal/api"
	"github.com/galchammat/kadeem/internal/database"
	"github.com/galchammat/kadeem/internal/logging"
	riot "github.com/galchammat/kadeem/internal/riot/api"
	"github.com/galchammat/kadeem/internal/service"
)

type Daemon struct {
	db      *database.DB
	matches *service.MatchService
	ranks   *service.RankService
}

func New(db *database.DB) *Daemon {
	riotClient := riot.NewClient()
	return &Daemon{
		db:      db,
		matches: service.NewMatchService(db, riotClient),
		ranks:   service.NewRankService(db, riotClient),
	}
}

func (d *Daemon) RunMatchSyncJob(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	d.syncMatches()

	for {
		select {
		case <-ctx.Done():
			logging.Info("Match sync job stopped")
			return
		case <-ticker.C:
			d.syncMatchesSafe()
		}
	}
}

func (d *Daemon) syncMatchesSafe() {
	defer func() {
		if r := recover(); r != nil {
			logging.Error("Panic in match sync job", "panic", r)
		}
	}()
	d.syncMatches()
}

func (d *Daemon) syncMatches() {
	logging.Info("Starting match sync")

	accounts, err := d.db.GetTrackedAccountsForSync()
	if err != nil {
		logging.Error("Failed to list accounts for sync", "error", err)
		return
	}

	for _, account := range accounts {
		if err := d.matches.SyncMatches(account); err != nil {
			logging.Error("Failed to sync matches for account",
				"puuid", account.PUUID, "error", err)
		} else {
			logging.Info("Synced matches for account", "puuid", account.PUUID)
		}
	}

	logging.Info("Match sync completed")
}

func (d *Daemon) RunRankSyncJob(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	d.syncRanks()

	for {
		select {
		case <-ctx.Done():
			logging.Info("Rank sync job stopped")
			return
		case <-ticker.C:
			d.syncRanksSafe()
		}
	}
}

func (d *Daemon) syncRanksSafe() {
	defer func() {
		if r := recover(); r != nil {
			logging.Error("Panic in rank sync job", "panic", r)
		}
	}()
	d.syncRanks()
}

func (d *Daemon) syncRanks() {
	logging.Info("Starting rank sync")

	accounts, err := d.db.GetTrackedAccountsForSync()
	if err != nil {
		logging.Error("Failed to list accounts for rank sync", "error", err)
		return
	}

	for _, account := range accounts {
		if err := d.ranks.SyncRank(&account); err != nil {
			logging.Error("Failed to sync rank for account",
				"puuid", account.PUUID, "error", err)
		} else {
			logging.Info("Synced rank for account", "puuid", account.PUUID)
		}
	}

	logging.Info("Rank sync completed")
}

// StartAPIServer starts the HTTP API server
func (d *Daemon) StartAPIServer(ctx context.Context) error {
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	logging.Info("Starting API server", "port", port)
	return api.StartServer(ctx, d.db, port)
}
