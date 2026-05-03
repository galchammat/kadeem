package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/galchammat/kadeem/internal/logging"
	platformdb "github.com/galchammat/kadeem/internal/platform/database"
	riot "github.com/galchammat/kadeem/internal/riot/api"
	riotstore "github.com/galchammat/kadeem/internal/riot/store"
	"github.com/galchammat/kadeem/internal/syncer"
	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

type localReplayWriter struct {
	root string
}

func (w localReplayWriter) WriteReplay(_ context.Context, key string, body io.Reader) error {
	path, err := w.path(key)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create artifact directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create replay artifact: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, body); err != nil {
		return fmt.Errorf("write replay artifact: %w", err)
	}

	return nil
}

func (w localReplayWriter) path(key string) (string, error) {
	cleanKey := filepath.Clean(key)
	if cleanKey == "." || strings.HasPrefix(cleanKey, "..") || filepath.IsAbs(cleanKey) {
		return "", fmt.Errorf("invalid artifact key %q", key)
	}

	return filepath.Join(w.root, cleanKey), nil
}

func main() {
	logging.Init(os.Stderr, slog.LevelInfo)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := platformdb.OpenDB()
	if err != nil {
		logging.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.SQL.Close()

	store := riotstore.New(db)
	accounts, err := store.GetTrackedAccountsForSync()
	if err != nil {
		logging.Error("failed to list accounts for replay sync", "error", err)
		os.Exit(1)
	}

	replayStore, err := riot.NewReplayStore(db.SQL)
	if err != nil {
		logging.Error("failed to create replay store", "error", err)
		os.Exit(1)
	}

	artifactRoot := os.Getenv("ARTIFACT_ROOT")
	if artifactRoot == "" {
		artifactRoot = "artifacts"
	}
	replayHandler, err := riot.NewReplayHandler(riot.NewClient(), localReplayWriter{root: artifactRoot}, accounts)
	if err != nil {
		logging.Error("failed to create replay handler", "error", err)
		os.Exit(1)
	}

	worker, err := syncer.NewArtifactWorker(syncer.ArtifactWorkerConfig{Logger: slog.Default()}, replayStore, replayHandler)
	if err != nil {
		logging.Error("failed to create artifact worker", "error", err)
		os.Exit(1)
	}

	logging.Info("starting lol artifact syncer", "artifact_root", artifactRoot)
	if err := worker.RunForever(ctx); err != nil {
		logging.Error("lol artifact syncer stopped", "error", err)
		os.Exit(1)
	}
}
