package syncer

import (
	"context"
	"fmt"
	"log/slog"
)

type Controller[E any] struct {
	source Source[E]
	logger *slog.Logger
}

type ControllerConfig struct {
	Logger *slog.Logger
}

func NewController[E any](cfg ControllerConfig, source Source[E]) (*Controller[E], error) {
	if source == nil {
		return nil, fmt.Errorf("source is nil")
	}

	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}

	return &Controller[E]{
		source: source,
		logger: cfg.Logger,
	}, nil
}

// Run performs metadata discovery and persistence.
func (c *Controller[E]) Run(ctx context.Context) error {
	if err := c.source.Sync(ctx); err != nil {
		return fmt.Errorf("sync source: %w", err)
	}

	c.logger.Info("synced source")
	return nil
}
