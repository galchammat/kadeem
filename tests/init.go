package tests

import (
	"os"

	"github.com/galchammat/kadeem/internal/logging"

	"github.com/joho/godotenv"
)

func init() {
	envFile := os.Getenv("ENV_FILE")
	logging.Info("Loading environment variables from .env file", "file", envFile)
	if err := godotenv.Load(envFile); err != nil {
		logging.Error("Failed to load .env file", "file", envFile, "error", err)
		panic(err)
	}
	os.Setenv("RUN_INTEGRATION_TESTS", "true") // This package only contains integration tests
}
