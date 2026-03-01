package api

import (
	"github.com/galchammat/kadeem/internal/api/middleware"
	"github.com/go-chi/chi/v5"
)

func (s *Server) setupRoutes(jwksURL string) {
	r := s.router

	// Middleware stack
	r.Use(middleware.RecoveryMiddleware)
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.CORSMiddleware(s.allowedOrigins))

	// Public endpoints (no auth required)
	r.Get("/health", s.healthHandler.Health)

	// API routes
	r.Route("/api/v0", func(r chi.Router) {
		// Public DataDragon endpoints
		r.Get("/datadragon/version", s.healthHandler.DataDragonVersion)
		r.Get("/datadragon/champions", s.dataDragonHandler.GetChampionData)
		r.Get("/datadragon/items", s.dataDragonHandler.GetItemData)
		r.Get("/datadragon/runes", s.dataDragonHandler.GetRuneData)
		r.Get("/datadragon/summoner-spells", s.dataDragonHandler.GetSummonerSpellData)

		// Protected routes (require authentication)
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(jwksURL))

			// Riot accounts
			r.Post("/riot/accounts", s.riotHandler.AddAccount)
			r.Get("/riot/accounts", s.riotHandler.ListAccounts)
			r.Get("/riot/accounts/{accountID}", s.riotHandler.GetAccount)
			r.Put("/riot/accounts/{accountID}", s.riotHandler.UpdateAccount)
			r.Delete("/riot/accounts/{accountID}", s.riotHandler.DeleteAccount)

			// Riot matches
			r.Post("/riot/accounts/{accountID}/matches/sync", s.riotHandler.SyncMatches)
			r.Get("/riot/matches", s.riotHandler.ListMatches)
			r.Post("/riot/matches/{matchID}/replay", s.riotHandler.SyncMatchReplay)
			r.Get("/riot/accounts/{accountID}/replays", s.riotHandler.FetchReplayURLs)
			r.Get("/riot/accounts/{accountID}/match-summaries", s.riotHandler.FetchMatchSummary)
			r.Post("/riot/matches/{matchID}/summary", s.riotHandler.SyncMatchSummary)

			// Riot ranks
			r.Get("/riot/accounts/{accountID}/rank-at-time", s.riotHandler.GetPlayerRankAtTime)
			r.Post("/riot/accounts/{accountID}/rank/sync", s.riotHandler.SyncRank)

			// Streamers
			r.Get("/streamers", s.livestreamHandler.ListStreamersWithDetails)
			r.Post("/streamers", s.livestreamHandler.AddStreamer)
			r.Delete("/streamers/{name}", s.livestreamHandler.DeleteStreamer)

			// Channels
			r.Post("/channels", s.livestreamHandler.AddChannel)
			r.Delete("/channels/{channelID}", s.livestreamHandler.DeleteChannel)

			// Broadcasts
			r.Post("/channels/{channelID}/broadcasts/sync", s.livestreamHandler.SyncBroadcasts)
			r.Get("/broadcasts", s.livestreamHandler.ListBroadcasts)

			// Stream events
			r.Post("/channels/{channelID}/events/sync", s.eventsHandler.SyncChannelEvents)
			r.Get("/channels/{channelID}/events", s.eventsHandler.ListChannelEvents)
			r.Get("/streamers/{streamerID}/events", s.eventsHandler.ListStreamerEvents)
		})
	})
}
