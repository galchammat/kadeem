package syncer

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type ArtifactWorker struct {
	store   ArtifactStore
	handler ArtifactHandler
	cfg     ArtifactWorkerConfig
}

type ArtifactWorkerConfig struct {
	Workers   int
	PollEvery time.Duration
	Logger    *slog.Logger
}

func NewArtifactWorker(
	cfg ArtifactWorkerConfig,
	store ArtifactStore,
	handler ArtifactHandler,
) (*ArtifactWorker, error) {
	if store == nil {
		return nil, fmt.Errorf("artifact store is nil")
	}
	if handler == nil {
		return nil, fmt.Errorf("artifact handler is nil")
	}

	if cfg.Workers <= 0 {
		cfg.Workers = 4
	}
	if cfg.PollEvery <= 0 {
		cfg.PollEvery = 10 * time.Second
	}
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}

	return &ArtifactWorker{
		store:   store,
		handler: handler,
		cfg:     cfg,
	}, nil
}

func (w *ArtifactWorker) RunForever(ctx context.Context) error {
	jobs := make(chan Artifact, w.cfg.Workers)

	for i := 0; i < w.cfg.Workers; i++ {
		go w.runWorker(ctx, jobs)
	}

	ticker := time.NewTicker(w.cfg.PollEvery)
	defer ticker.Stop()

	for {
		limit := cap(jobs) - len(jobs)
		if limit <= 0 {
			select {
			case <-ticker.C:
			case <-ctx.Done():
				return nil
			}
			continue
		}

		artifacts, err := w.store.ClaimPending(ctx, limit)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			w.cfg.Logger.Error("claim artifacts", "error", err)
		}

		for _, artifact := range artifacts {
			select {
			case jobs <- artifact:
			case <-ctx.Done():
				return nil
			}
		}

		select {
		case <-ticker.C:
		case <-ctx.Done():
			return nil
		}
	}
}

// ProcessPending claims and processes one batch of pending artifacts.
func (w *ArtifactWorker) ProcessPending(ctx context.Context) (int, error) {
	artifacts, err := w.store.ClaimPending(ctx, w.cfg.Workers)
	if err != nil {
		return 0, fmt.Errorf("claim artifacts: %w", err)
	}

	for _, artifact := range artifacts {
		if err := ctx.Err(); err != nil {
			return len(artifacts), err
		}
		w.process(ctx, artifact)
	}

	return len(artifacts), nil
}

func (w *ArtifactWorker) runWorker(ctx context.Context, jobs <-chan Artifact) {
	for {
		select {
		case artifact, ok := <-jobs:
			if !ok {
				return
			}
			w.process(ctx, artifact)
		case <-ctx.Done():
			return
		}
	}
}

func (w *ArtifactWorker) process(ctx context.Context, artifact Artifact) {
	s3Key, err := w.handler.Process(ctx, artifact)
	if err != nil {
		if markErr := w.store.MarkFailed(ctx, artifact.ID, err); markErr != nil {
			w.cfg.Logger.Error("mark artifact failed", "id", artifact.ID, "process_error", err, "error", markErr)
		}
		return
	}

	if err := w.store.MarkDone(ctx, artifact.ID, s3Key); err != nil {
		w.cfg.Logger.Error("mark artifact done", "id", artifact.ID, "error", err)
	}
}
