package main

import (
	"context"
	"embed"

	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/galchammat/kadeem/internal/database"
	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/riot"
)

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		logging.Warn("No .env file found or unable to load")
	}
}

func main() {
	ctx := context.Background()
	app := NewApp()
	DB, err := database.OpenDB()
	if err != nil {
		logging.Error("Failed to open database: %v", err)
		return
	}
	riotHandler := riot.NewRiotClient(ctx, DB)

	err = wails.Run(&options.App{
		Title:       "kadeem",
		Width:       1024,
		Height:      768,
		StartHidden: true, // start without showing the window
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(wailsCtx context.Context) {
			// You can still use Wails context if needed for Wails-specific features
			app.startup(wailsCtx)

		},
		Bind: []interface{}{
			app,
			riotHandler,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
