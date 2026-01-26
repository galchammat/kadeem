package daemon

import (
	"context"
	"time"

	"github.com/galchammat/kadeem/internal/database"
	"github.com/galchammat/kadeem/internal/logging"
	riot "github.com/galchammat/kadeem/internal/riot/api"
)

type Daemon struct {
	db         *database.DB
	riotClient *riot.RiotClient
}

func New(db *database.DB) *Daemon {
	return &Daemon{
		db:         db,
		riotClient: riot.NewRiotClient(context.Background(), db),
	}
}

func (d *Daemon) RunMatchSyncJob(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Run immediately on start
	d.syncMatches()

	for {
		select {
		case <-ctx.Done():
			logging.Info("Match sync job stopped")
			return
		case <-ticker.C:
			d.syncMatches()
		}
	}
}

func (d *Daemon) syncMatches() {
	logging.Info("Starting match sync")

	// Get all accounts
	accounts, err := d.db.ListRiotAccounts(nil)
	if err != nil {
		logging.Error("Failed to list accounts for sync", "error", err)
		return
	}

	// Sync matches for each account
	for _, account := range accounts {
		if err := d.riotClient.SyncMatches(account); err != nil {
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
			d.syncRanks()
		}
	}
}

func (d *Daemon) syncRanks() {
	logging.Info("Starting rank sync")

	accounts, err := d.db.ListRiotAccounts(nil)
	if err != nil {
		logging.Error("Failed to list accounts for rank sync", "error", err)
		return
	}

	for _, account := range accounts {
		if err := d.riotClient.SyncRank(&account); err != nil {
			logging.Error("Failed to sync rank for account",
				"puuid", account.PUUID, "error", err)
		} else {
			logging.Info("Synced rank for account", "puuid", account.PUUID)
		}
	}

	logging.Info("Rank sync completed")
}
