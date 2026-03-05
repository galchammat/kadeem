package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/service"
	"github.com/go-chi/chi/v5"
)

// EventsHandler handles stream event HTTP requests.
type EventsHandler struct {
	events *service.StreamEventsService
}

// NewEventsHandler creates a new EventsHandler.
func NewEventsHandler(events *service.StreamEventsService) *EventsHandler {
	return &EventsHandler{events: events}
}

// SyncChannelEvents triggers a sync of hype train and clip events for the given channel.
func (h *EventsHandler) SyncChannelEvents(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelID")
	if err := h.events.SyncChannelEvents(channelID); err != nil {
		logging.Error("failed to sync channel events", "channel_id", channelID, "error", err)
		respondError(w, http.StatusInternalServerError, "failed to sync events")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListChannelEvents returns stream events for a specific channel.
func (h *EventsHandler) ListChannelEvents(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelID")
	from, to, limit, offset := parseEventParams(r)

	events, err := h.events.ListChannelEvents(channelID, from, to, limit, offset)
	if err != nil {
		logging.Error("failed to list channel events", "channel_id", channelID, "error", err)
		respondError(w, http.StatusInternalServerError, "failed to list events")
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"events": events,
		"count":  len(events),
	})
}

// ListStreamerEvents returns stream events for all channels of a streamer.
func (h *EventsHandler) ListStreamerEvents(w http.ResponseWriter, r *http.Request) {
	streamerID, err := strconv.ParseInt(chi.URLParam(r, "streamerID"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid streamer ID")
		return
	}

	from, to, limit, offset := parseEventParams(r)

	events, err := h.events.ListStreamerEvents(streamerID, from, to, limit, offset)
	if err != nil {
		logging.Error("failed to list streamer events", "streamer_id", streamerID, "error", err)
		respondError(w, http.StatusInternalServerError, "failed to list events")
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"events": events,
		"count":  len(events),
	})
}

// parseEventParams extracts from, to, limit, and offset from query params with sensible defaults.
func parseEventParams(r *http.Request) (from, to int64, limit, offset int) {
	q := r.URL.Query()

	from, _ = strconv.ParseInt(q.Get("from"), 10, 64)

	to, _ = strconv.ParseInt(q.Get("to"), 10, 64)
	if to == 0 {
		to = time.Now().Unix()
	}

	limit, _ = strconv.Atoi(q.Get("limit"))
	if limit <= 0 || limit > 500 {
		limit = 100
	}

	offset, _ = strconv.Atoi(q.Get("offset"))
	if offset < 0 {
		offset = 0
	}

	return from, to, limit, offset
}
