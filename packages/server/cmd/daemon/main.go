package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/galchammat/kadeem/internal/daemon"
	"github.com/galchammat/kadeem/internal/database"
	"github.com/galchammat/kadeem/internal/logging"
)

func main() {
	logging.Init(os.Stderr, slog.LevelInfo)
	logging.Info("Starting Kadeem daemon")

	db, err := database.OpenDB()
	if err != nil {
		logging.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.SQL.Close()

	d := daemon.New(db)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		d.RunMatchSyncJob(ctx, 15*time.Minute)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		d.RunRankSyncJob(ctx, 15*time.Minute)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := d.StartAPIServer(ctx); err != nil {
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
