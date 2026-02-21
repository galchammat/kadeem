package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/galchammat/kadeem/internal/api"
	"github.com/galchammat/kadeem/internal/database"
	"github.com/galchammat/kadeem/internal/logging"
	riot "github.com/galchammat/kadeem/internal/riot/api"
	"github.com/galchammat/kadeem/internal/service"
	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

type daemon struct {
	db      *database.DB
	matches *service.MatchService
	ranks   *service.RankService
}

func main() {
	logging.Init(os.Stderr, slog.LevelInfo)
	logging.Info("Starting Kadeem daemon")

	db, err := database.OpenDB()
	if err != nil {
		logging.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.SQL.Close()

	riotClient := riot.NewClient()
	d := &daemon{
		db:      db,
		matches: service.NewMatchService(db, riotClient),
		ranks:   service.NewRankService(db, riotClient),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		d.runSyncLoop(ctx, 15*time.Minute, "match", d.syncMatches)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		d.runSyncLoop(ctx, 15*time.Minute, "rank", d.syncRanks)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		port := os.Getenv("API_PORT")
		if port == "" {
			port = "8080"
		}
		logging.Info("Starting API server", "port", port)
		if err := api.StartServer(ctx, db, port); err != nil {
			logging.Error("API server stopped", "error", err)
		}
	}()

	logging.Info("Daemon running (API + background jobs), press Ctrl+C to stop")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logging.Info("Shutting down daemon")
	cancel()
	wg.Wait()
	logging.Info("Daemon stopped")
}

func (d *daemon) runSyncLoop(ctx context.Context, interval time.Duration, name string, fn func()) {
	safeRun := func() {
		defer func() {
			if r := recover(); r != nil {
				logging.Error("Panic in sync job", "job", name, "panic", r)
			}
		}()
		fn()
	}

	safeRun()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			logging.Info("Sync job stopped", "job", name)
			return
		case <-ticker.C:
			safeRun()
		}
	}
}

func (d *daemon) syncMatches() {
	logging.Info("Starting match sync")
	accounts, err := d.db.GetTrackedAccountsForSync()
	if err != nil {
		logging.Error("Failed to list accounts for sync", "error", err)
		return
	}
	for _, account := range accounts {
		if err := d.matches.SyncMatches(account); err != nil {
			logging.Error("Failed to sync matches", "puuid", account.PUUID, "error", err)
		} else {
			logging.Info("Synced matches", "puuid", account.PUUID)
		}
	}
	logging.Info("Match sync completed")
}

func (d *daemon) syncRanks() {
	logging.Info("Starting rank sync")
	accounts, err := d.db.GetTrackedAccountsForSync()
	if err != nil {
		logging.Error("Failed to list accounts for rank sync", "error", err)
		return
	}
	for _, account := range accounts {
		if err := d.ranks.SyncRank(&account); err != nil {
			logging.Error("Failed to sync rank", "puuid", account.PUUID, "error", err)
		} else {
			logging.Info("Synced rank", "puuid", account.PUUID)
		}
	}
	logging.Info("Rank sync completed")
}
