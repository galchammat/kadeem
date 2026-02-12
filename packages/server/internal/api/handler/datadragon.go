package handler

import (
	"encoding/json"
	"net/http"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/riot/datadragon"
)

type DataDragonHandler struct {
	client *datadragon.DataDragonClient
}

func NewDataDragonHandler(client *datadragon.DataDragonClient) *DataDragonHandler {
	return &DataDragonHandler{client: client}
}

func (h *DataDragonHandler) GetChampionData(w http.ResponseWriter, r *http.Request) {
	locale := r.URL.Query().Get("locale")
	if locale == "" {
		locale = "en_US"
	}

	data, err := h.client.GetChampionData(locale)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, data)
}

func (h *DataDragonHandler) GetItemData(w http.ResponseWriter, r *http.Request) {
	locale := r.URL.Query().Get("locale")
	if locale == "" {
		locale = "en_US"
	}

	data, err := h.client.GetItemData(locale)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, data)
}

func (h *DataDragonHandler) GetRuneData(w http.ResponseWriter, r *http.Request) {
	locale := r.URL.Query().Get("locale")
	if locale == "" {
		locale = "en_US"
	}

	data, err := h.client.GetRuneData(locale)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, data)
}

func (h *DataDragonHandler) GetSummonerSpellData(w http.ResponseWriter, r *http.Request) {
	locale := r.URL.Query().Get("locale")
	if locale == "" {
		locale = "en_US"
	}

	data, err := h.client.GetSummonerSpellData(locale)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, data)
}

// Helper functions
func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logging.Error("Failed to encode JSON response", "status", status, "error", err)
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		logging.Error("Failed to encode error response", "status", status, "error", err)
	}
}
