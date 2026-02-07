package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/galchammat/kadeem/internal/api/middleware"
	apiModels "github.com/galchammat/kadeem/internal/api/models"
	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/model"
	"github.com/galchammat/kadeem/internal/service"
	"github.com/go-chi/chi/v5"
)

type LivestreamHandler struct {
	streamers *service.StreamerService
}

func NewLivestreamHandler(streamers *service.StreamerService) *LivestreamHandler {
	return &LivestreamHandler{streamers: streamers}
}

// ListStreamersWithDetails returns all streamers with details
func (h *LivestreamHandler) ListStreamersWithDetails(w http.ResponseWriter, r *http.Request) {
	streamers, err := h.streamers.ListStreamersWithDetails()
	if err != nil {
		logging.Error("Failed to list streamers", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to list streamers")
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"streamers": streamers,
		"count":     len(streamers),
	})
}

// AddStreamer creates a streamer
func (h *LivestreamHandler) AddStreamer(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)

	var req apiModels.AddStreamerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	id, err := h.streamers.AddStreamer(req.Name)
	if err != nil {
		logging.Error("Failed to add streamer", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to add streamer")
		return
	}

	// TODO: Track streamer for user when tracking is implemented
	_ = userID
	_ = id

	respondJSON(w, http.StatusOK, apiModels.SuccessResponse{
		Message: "Streamer added successfully",
	})
}

// DeleteStreamer deletes a streamer
func (h *LivestreamHandler) DeleteStreamer(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	deleted, err := h.streamers.DeleteStreamer(name)
	if err != nil {
		logging.Error("Failed to delete streamer", "name", name, "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to delete streamer")
		return
	}

	if !deleted {
		respondError(w, http.StatusNotFound, "Streamer not found")
		return
	}

	respondJSON(w, http.StatusOK, apiModels.SuccessResponse{
		Message: "Streamer deleted successfully",
	})
}

// AddChannel adds a channel to a streamer
func (h *LivestreamHandler) AddChannel(w http.ResponseWriter, r *http.Request) {
	var req apiModels.AddChannelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	channel := model.Channel{
		StreamerID:  int64(req.StreamerID),
		ChannelName: req.ChannelName,
		ID:          req.ChannelID,
		Platform:    req.Platform,
	}

	saved, err := h.streamers.AddChannel(channel)
	if err != nil {
		logging.Error("Failed to add channel", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to add channel")
		return
	}

	if !saved {
		respondError(w, http.StatusConflict, "Channel already exists")
		return
	}

	respondJSON(w, http.StatusCreated, apiModels.SuccessResponse{
		Message: "Channel added successfully",
	})
}

// DeleteChannel deletes a channel
func (h *LivestreamHandler) DeleteChannel(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	channelID := chi.URLParam(r, "channelID")

	if !middleware.IsAdmin(r) {
		// TODO: Check channel ownership through streamer
		_ = userID
	}

	deleted, err := h.streamers.DeleteChannel(channelID)
	if err != nil {
		logging.Error("Failed to delete channel", "channelID", channelID, "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to delete channel")
		return
	}

	if !deleted {
		respondError(w, http.StatusNotFound, "Channel not found")
		return
	}

	respondJSON(w, http.StatusOK, apiModels.SuccessResponse{
		Message: "Channel deleted successfully",
	})
}

// SyncBroadcasts syncs broadcasts for a channel
func (h *LivestreamHandler) SyncBroadcasts(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	channelID := chi.URLParam(r, "channelID")

	if !middleware.IsAdmin(r) {
		// TODO: Check channel ownership through streamer
		_ = userID
	}

	channel := model.Channel{ID: channelID}

	err := h.streamers.SyncBroadcasts(channel)
	if err != nil {
		logging.Error("Failed to sync broadcasts", "channelID", channelID, "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to sync broadcasts")
		return
	}

	respondJSON(w, http.StatusOK, apiModels.SuccessResponse{
		Message: "Broadcasts synced successfully",
	})
}

// ListBroadcasts lists broadcasts for a channel
func (h *LivestreamHandler) ListBroadcasts(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)

	channelID := r.URL.Query().Get("channelID")
	if channelID == "" {
		respondError(w, http.StatusBadRequest, "channelID is required")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit == 0 {
		limit = 20
	}

	if !middleware.IsAdmin(r) {
		// TODO: Check channel ownership through streamer
		_ = userID
	}

	filter := &model.Broadcast{ChannelID: channelID}
	broadcasts, err := h.streamers.ListBroadcasts(filter, limit, offset)
	if err != nil {
		logging.Error("Failed to list broadcasts", "channelID", channelID, "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to list broadcasts")
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"broadcasts": broadcasts,
		"count":      len(broadcasts),
	})
}
