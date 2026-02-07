package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/galchammat/kadeem/internal/api/handler"
	"github.com/galchammat/kadeem/internal/database"
	"github.com/galchammat/kadeem/internal/logging"
	riot "github.com/galchammat/kadeem/internal/riot/api"
	"github.com/galchammat/kadeem/internal/riot/datadragon"
	"github.com/galchammat/kadeem/internal/service"
	"github.com/galchammat/kadeem/internal/twitch"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	router            *chi.Mux
	httpServer        *http.Server
	allowedOrigins    []string
	healthHandler     *handler.HealthHandler
	riotHandler       *handler.RiotHandler
	dataDragonHandler *handler.DataDragonHandler
	livestreamHandler *handler.LivestreamHandler
}

// NewServer creates a new API server
func NewServer(db *database.DB, port string) *Server {
	// Create clients
	ctx := context.Background()
	riotClient := riot.NewClient()
	dataDragonClient := datadragon.NewDataDragonClient(ctx, "bin/datadragon")
	twitchClient := twitch.NewTwitchClient(ctx)

	// Create services
	accountSvc := service.NewAccountService(db, riotClient)
	matchSvc := service.NewMatchService(db, riotClient)
	rankSvc := service.NewRankService(db, riotClient)
	streamerSvc := service.NewStreamerService(db, twitchClient)

	// Get frontend domain from env
	frontendDomain := os.Getenv("FRONTEND_DOMAIN")
	if frontendDomain == "" {
		frontendDomain = "cyanlab.cc"
	}

	allowedOrigins := []string{"https://" + frontendDomain, "http://localhost:5173"}

	version := os.Getenv("API_VERSION")
	if version == "" {
		version = "dev"
	}

	s := &Server{
		router:            chi.NewRouter(),
		allowedOrigins:    allowedOrigins,
		healthHandler:     handler.NewHealthHandler(version, db, dataDragonClient),
		riotHandler:       handler.NewRiotHandler(db, accountSvc, matchSvc, rankSvc),
		dataDragonHandler: handler.NewDataDragonHandler(dataDragonClient),
		livestreamHandler: handler.NewLivestreamHandler(streamerSvc),
	}

	// Setup routes
	s.setupRoutes()

	// Create HTTP server
	s.httpServer = &http.Server{
		Addr:              "127.0.0.1:" + port,
		Handler:           s.router,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	return s
}

// Start starts the HTTP server
func (s *Server) Start() error {
	logging.Info("Starting API server", "addr", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	logging.Info("Shutting down API server")
	return s.httpServer.Shutdown(ctx)
}

// StartServer starts the API server (for daemon integration)
func StartServer(ctx context.Context, db *database.DB, port string) error {
	server := NewServer(db, port)

	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			logging.Error("API server error", "error", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("failed to shutdown API server: %w", err)
	}

	return nil
}
