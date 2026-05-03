package syncer

import (
	"context"
	"fmt"
	"log/slog"
)

type MetadataSyncer[E any] struct {
	source Source[E]
	logger *slog.Logger
}

type MetadataSyncerConfig struct {
	Logger *slog.Logger
}

func NewMetadataSyncer[E any](cfg MetadataSyncerConfig, source Source[E]) (*MetadataSyncer[E], error) {
	if source == nil {
		return nil, fmt.Errorf("source is nil")
	}
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}

	return &MetadataSyncer[E]{
		source: source,
		logger: cfg.Logger,
	}, nil
}

// RunOnce performs metadata discovery and persistence once.
func (s *MetadataSyncer[E]) RunOnce(ctx context.Context) error {
	if err := s.source.Sync(ctx); err != nil {
		return fmt.Errorf("sync source: %w", err)
	}
	s.logger.Info("synced metadata")

	return nil
}

// Run performs metadata discovery and persistence once.
func (s *MetadataSyncer[E]) Run(ctx context.Context) error {
	return s.RunOnce(ctx)
}
