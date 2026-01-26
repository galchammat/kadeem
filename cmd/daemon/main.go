package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/galchammat/kadeem/internal/daemon"
	"github.com/galchammat/kadeem/internal/database"
	"github.com/galchammat/kadeem/internal/logging"
)

func main() {
	logging.Info("Starting Kadeem daemon")

	// Connect to PostgreSQL
	db, err := database.OpenDB()
	if err != nil {
		logging.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.SQL.Close()

	// Create daemon
	d := daemon.New(db)

	// Start background jobs
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go d.RunMatchSyncJob(ctx, 15*time.Minute)
	go d.RunRankSyncJob(ctx, 15*time.Minute)

	logging.Info("Daemon running, press Ctrl+C to stop")

	// Wait for interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logging.Info("Shutting down daemon")
}
