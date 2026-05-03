package main

import (
	"log/slog"
	"os"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	logging.Init(os.Stderr, slog.LevelInfo)
	logging.Error("twitch artifact sync is not implemented yet")
	os.Exit(1)
}
