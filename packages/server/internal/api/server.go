package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/galchammat/kadeem/internal/api/handler"
	"github.com/galchammat/kadeem/internal/logging"
	platformdb "github.com/galchammat/kadeem/internal/platform/database"
	riotapi "github.com/galchammat/kadeem/internal/riot/api"
	"github.com/galchammat/kadeem/internal/riot/datadragon"
	riotpostgres "github.com/galchammat/kadeem/internal/riot/postgres"
	"github.com/galchammat/kadeem/internal/service"
	twitchapi "github.com/galchammat/kadeem/internal/twitch/api"
	twitchstore "github.com/galchammat/kadeem/internal/twitch/store"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	router            *chi.Mux
	httpServer        *http.Server
	allowedOrigins    []string
	jwksURL           string
	healthHandler     *handler.HealthHandler
	riotHandler       *handler.RiotHandler
	dataDragonHandler *handler.DataDragonHandler
	livestreamHandler *handler.LivestreamHandler
	eventsHandler     *handler.EventsHandler
}

// NewServer creates a new API server
func NewServer(db *platformdb.DB, riotStore *riotpostgres.DB, twitchStore *twitchstore.Store, port string) *Server {
	// Create clients
	ctx := context.Background()
	riotClient := riotapi.NewClient()
	dataDragonClient := datadragon.NewDataDragonClient(ctx, "bin/datadragon")
	twitchClient := twitchapi.NewTwitchClient(ctx)

	// Create services
	accountSvc := service.NewAccountService(riotStore, riotClient)
	matchSvc := service.NewMatchService(riotStore, riotClient)
	rankSvc := service.NewRankService(riotStore, riotClient)
	streamerSvc := service.NewStreamerService(twitchStore, twitchClient)
	streamEventsSvc := service.NewStreamEventsService(twitchStore, twitchClient)

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

	jwksURL := os.Getenv("SUPABASE_JWKS_URL")
	if jwksURL == "" {
		// Default to the known Supabase JWKS URL
		jwksURL = "https://seijlvqsunpbzwuydvze.supabase.co/auth/v1/.well-known/jwks.json"
		logging.Info("SUPABASE_JWKS_URL not set, using default")
	}

	s := &Server{
		router:            chi.NewRouter(),
		allowedOrigins:    allowedOrigins,
		jwksURL:           jwksURL,
		healthHandler:     handler.NewHealthHandler(version, db, dataDragonClient),
		riotHandler:       handler.NewRiotHandler(riotStore, twitchStore, accountSvc, matchSvc, rankSvc),
		dataDragonHandler: handler.NewDataDragonHandler(dataDragonClient),
		livestreamHandler: handler.NewLivestreamHandler(streamerSvc),
		eventsHandler:     handler.NewEventsHandler(streamEventsSvc),
	}

	// Setup routes
	s.setupRoutes(s.jwksURL)

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
func StartServer(ctx context.Context, db *platformdb.DB, riotStore *riotpostgres.DB, twitchStore *twitchstore.Store, port string) error {
	server := NewServer(db, riotStore, twitchStore, port)

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
