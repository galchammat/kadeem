package handler

import (
	"encoding/json"
	"net/http"

	"github.com/galchammat/kadeem/internal/api/models"
	"github.com/galchammat/kadeem/internal/database"
	"github.com/galchammat/kadeem/internal/riot/datadragon"
)

type HealthHandler struct {
	version          string
	db               *database.DB
	dataDragonClient *datadragon.DataDragonClient
}

func NewHealthHandler(version string, db *database.DB, ddClient *datadragon.DataDragonClient) *HealthHandler {
	return &HealthHandler{
		version:          version,
		db:               db,
		dataDragonClient: ddClient,
	}
}

// Health returns API health status
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	status := "ok"
	if err := h.db.SQL.PingContext(r.Context()); err != nil {
		status = "degraded"
	}

	response := models.HealthResponse{
		Status:  status,
		Version: h.version,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DataDragonVersion returns current DataDragon version
func (h *HealthHandler) DataDragonVersion(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"version": h.dataDragonClient.GetVersion(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
