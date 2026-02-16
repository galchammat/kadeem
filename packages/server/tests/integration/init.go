package tests

import (
	"os"

	"github.com/galchammat/kadeem/internal/logging"

	"github.com/joho/godotenv"
)

func init() {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		return
	}
	logging.Info("Loading environment variables from .env file", "file", ".env")
	if err := godotenv.Load(); err != nil {
		logging.Error("Failed to load .env file", "file", ".env", "error", err)
		return
	}
	os.Setenv("RUN_INTEGRATION_TESTS", "true")
}
