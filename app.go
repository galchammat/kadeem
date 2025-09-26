package main

import (
	"context"
	"fmt"

	"clipdeem/internal/riot"
)

// App struct
type App struct {
	ctx         context.Context
	RiotHandler *riot.RiotHandler
	// TwitchHandler twitch.TwitchHandler
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		// TwitchHandler: twitch.NewTwitchHandler(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.RiotHandler = riot.NewRiotHandler(ctx)
	// a.TwitchHandler.SetContext(ctx)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time brotha!", name)
}
