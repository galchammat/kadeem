package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/galchammat/kadeem/internal/api/middleware"
	apiModels "github.com/galchammat/kadeem/internal/api/models"
	"github.com/galchammat/kadeem/internal/database"
	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/model"
	"github.com/galchammat/kadeem/internal/service"
	"github.com/go-chi/chi/v5"
)

type RiotHandler struct {
	db       *database.DB
	accounts *service.AccountService
	matches  *service.MatchService
	ranks    *service.RankService
}

func NewRiotHandler(db *database.DB, accounts *service.AccountService, matches *service.MatchService, ranks *service.RankService) *RiotHandler {
	return &RiotHandler{
		db:       db,
		accounts: accounts,
		matches:  matches,
		ranks:    ranks,
	}
}

// AddAccount finds or creates an account and tracks it for the user (idempotent)
func (h *RiotHandler) AddAccount(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)

	var req apiModels.AddAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.StreamerID <= 0 {
		respondError(w, http.StatusBadRequest, "Missing streamer_id")
		return
	}

	streamer, err := h.db.GetStreamerByID(req.StreamerID)
	if err != nil {
		logging.Error("Failed to look up streamer", "streamerID", req.StreamerID, "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to add account")
		return
	}
	if streamer == nil {
		respondError(w, http.StatusBadRequest, "Invalid streamer_id")
		return
	}

	account, err := h.db.FindOrCreateRiotAccount(req.GameName, req.TagLine, req.Region, req.StreamerID)
	if err != nil {
		logging.Error("Failed to find or create account", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to add account")
		return
	}

	err = h.db.TrackAccount(userID, account.PUUID)
	if err != nil {
		logging.Error("Failed to track account", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to track account")
		return
	}

	respondJSON(w, http.StatusOK, apiModels.SuccessResponse{
		Message: "Account tracked successfully",
	})
}

// ListAccounts returns accounts this user is tracking
func (h *RiotHandler) ListAccounts(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)

	accounts, err := h.db.ListTrackedAccounts(userID)
	if err != nil {
		logging.Error("Failed to list tracked accounts", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to list accounts")
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"accounts": accounts,
		"count":    len(accounts),
	})
}

// GetAccount returns a specific account by PUUID (with tracking check)
func (h *RiotHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	accountPUUID := chi.URLParam(r, "accountID")
	if accountPUUID == "" {
		respondError(w, http.StatusBadRequest, "Missing account ID")
		return
	}

	isTracking, err := h.db.IsTrackingAccount(userID, accountPUUID)
	if err != nil {
		logging.Error("Failed to check tracking", "error", err)
		respondError(w, http.StatusInternalServerError, "Internal error")
		return
	}
	if !isTracking {
		respondError(w, http.StatusForbidden, "Not tracking this account")
		return
	}

	account, err := h.db.GetRiotAccount(accountPUUID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Account not found")
		return
	}

	respondJSON(w, http.StatusOK, account)
}

// UpdateAccount updates an existing account
func (h *RiotHandler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	accountPUUID := chi.URLParam(r, "accountID")
	if accountPUUID == "" {
		respondError(w, http.StatusBadRequest, "Missing account ID")
		return
	}

	var req apiModels.UpdateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	isTracking, err := h.db.IsTrackingAccount(userID, accountPUUID)
	if err != nil || !isTracking {
		respondError(w, http.StatusForbidden, "Not tracking this account")
		return
	}

	account, err := h.db.GetRiotAccount(accountPUUID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Account not found")
		return
	}

	err = h.accounts.UpdateAccount(req.Region, req.GameName, req.TagLine, account.PUUID)
	if err != nil {
		logging.Error("Failed to update account", "puuid", accountPUUID, "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to update account")
		return
	}

	respondJSON(w, http.StatusOK, apiModels.SuccessResponse{
		Message: "Account updated successfully",
	})
}

// DeleteAccount untracks an account for the user
func (h *RiotHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	accountPUUID := chi.URLParam(r, "accountID")
	if accountPUUID == "" {
		respondError(w, http.StatusBadRequest, "Missing account ID")
		return
	}

	err := h.db.UntrackAccount(userID, accountPUUID)
	if err != nil {
		logging.Error("Failed to untrack account", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to delete account")
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
}

// SyncMatches syncs matches for an account
func (h *RiotHandler) SyncMatches(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	accountPUUID := chi.URLParam(r, "accountID")
	if accountPUUID == "" {
		respondError(w, http.StatusBadRequest, "Missing account ID")
		return
	}

	isTracking, err := h.db.IsTrackingAccount(userID, accountPUUID)
	if err != nil || !isTracking {
		respondError(w, http.StatusForbidden, "Not tracking this account")
		return
	}

	account, err := h.db.GetRiotAccount(accountPUUID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Account not found")
		return
	}

	err = h.matches.SyncMatches(*account)
	if err != nil {
		logging.Error("Failed to sync matches", "puuid", accountPUUID, "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to sync matches")
		return
	}

	respondJSON(w, http.StatusOK, apiModels.SuccessResponse{
		Message: "Matches synced successfully",
	})
}

// ListMatches returns matches with filters
func (h *RiotHandler) ListMatches(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)

	puuid := r.URL.Query().Get("puuid")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit == 0 {
		limit = 20
	}

	var account *model.LeagueOfLegendsAccount
	if puuid != "" {
		acc, err := h.db.GetRiotAccount(puuid)
		if err != nil {
			respondError(w, http.StatusNotFound, "Account not found")
			return
		}

		isTracking, err := h.db.IsTrackingAccount(userID, acc.PUUID)
		if err != nil || !isTracking {
			respondError(w, http.StatusForbidden, "Not tracking this account")
			return
		}
		account = acc
	}

	filter := &model.LolMatchFilter{}
	if puuid != "" {
		filter.PUUID = &puuid
	}

	matches, err := h.matches.ListMatches(filter, account, limit, offset)
	if err != nil {
		logging.Error("Failed to list matches", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to list matches")
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"matches": matches,
		"count":   len(matches),
	})
}

// SyncMatchReplay syncs replay for a match
func (h *RiotHandler) SyncMatchReplay(w http.ResponseWriter, r *http.Request) {
	matchIDStr := chi.URLParam(r, "matchID")
	matchID, err := strconv.ParseInt(matchIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid match ID")
		return
	}

	var req apiModels.SyncMatchReplayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = h.matches.SyncMatchReplay(matchID, req.URL)
	if err != nil {
		logging.Error("Failed to sync replay", "matchID", matchID, "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to sync replay")
		return
	}

	respondJSON(w, http.StatusOK, apiModels.SuccessResponse{
		Message: "Replay synced successfully",
	})
}

// FetchReplayURLs fetches replay URLs for an account
func (h *RiotHandler) FetchReplayURLs(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	accountPUUID := chi.URLParam(r, "accountID")
	if accountPUUID == "" {
		respondError(w, http.StatusBadRequest, "Missing account ID")
		return
	}

	isTracking, err := h.db.IsTrackingAccount(userID, accountPUUID)
	if err != nil || !isTracking {
		respondError(w, http.StatusForbidden, "Not tracking this account")
		return
	}

	account, err := h.db.GetRiotAccount(accountPUUID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Account not found")
		return
	}

	region := r.URL.Query().Get("region")
	if region == "" {
		region = account.Region
	}

	urls, err := h.matches.FetchReplayURLs(account.PUUID, region)
	if err != nil {
		logging.Error("Failed to fetch replay URLs", "puuid", accountPUUID, "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to fetch replay URLs")
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"urls": urls,
	})
}

// FetchMatchSummary fetches match IDs for an account
func (h *RiotHandler) FetchMatchSummary(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	accountPUUID := chi.URLParam(r, "accountID")
	if accountPUUID == "" {
		respondError(w, http.StatusBadRequest, "Missing account ID")
		return
	}

	isTracking, err := h.db.IsTrackingAccount(userID, accountPUUID)
	if err != nil || !isTracking {
		respondError(w, http.StatusForbidden, "Not tracking this account")
		return
	}

	account, err := h.db.GetRiotAccount(accountPUUID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Account not found")
		return
	}

	matchIDs, err := h.matches.FetchMatchIDs(account.PUUID, account.Region, nil)
	if err != nil {
		logging.Error("Failed to fetch match IDs", "puuid", accountPUUID, "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to fetch match summaries")
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"match_ids": matchIDs,
	})
}

// SyncMatchSummary syncs match summary
func (h *RiotHandler) SyncMatchSummary(w http.ResponseWriter, r *http.Request) {
	matchIDStr := chi.URLParam(r, "matchID")
	matchID, err := strconv.ParseInt(matchIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid match ID")
		return
	}

	var req apiModels.SyncMatchSummaryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = h.matches.SyncMatchSummary(matchID, req.FullMatchID, req.Region)
	if err != nil {
		logging.Error("Failed to sync match summary", "matchID", matchID, "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to sync match summary")
		return
	}

	respondJSON(w, http.StatusOK, apiModels.SuccessResponse{
		Message: "Match summary synced successfully",
	})
}

// GetPlayerRankAtTime gets rank at specific time
func (h *RiotHandler) GetPlayerRankAtTime(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	accountPUUID := chi.URLParam(r, "accountID")
	if accountPUUID == "" {
		respondError(w, http.StatusBadRequest, "Missing account ID")
		return
	}

	queueID, err := strconv.Atoi(r.URL.Query().Get("queueID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid queueID")
		return
	}

	timestamp, err := strconv.ParseInt(r.URL.Query().Get("timestamp"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid timestamp")
		return
	}

	isTracking, err := h.db.IsTrackingAccount(userID, accountPUUID)
	if err != nil || !isTracking {
		respondError(w, http.StatusForbidden, "Not tracking this account")
		return
	}

	account, err := h.db.GetRiotAccount(accountPUUID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Account not found")
		return
	}

	rank, err := h.accounts.GetPlayerRankAtTime(account.PUUID, queueID, timestamp)
	if err != nil {
		logging.Error("Failed to get rank at time", "puuid", accountPUUID, "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to get rank")
		return
	}

	respondJSON(w, http.StatusOK, rank)
}

// SyncRank syncs rank for an account
func (h *RiotHandler) SyncRank(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	accountPUUID := chi.URLParam(r, "accountID")
	if accountPUUID == "" {
		respondError(w, http.StatusBadRequest, "Missing account ID")
		return
	}

	isTracking, err := h.db.IsTrackingAccount(userID, accountPUUID)
	if err != nil || !isTracking {
		respondError(w, http.StatusForbidden, "Not tracking this account")
		return
	}

	account, err := h.db.GetRiotAccount(accountPUUID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Account not found")
		return
	}

	err = h.ranks.SyncRank(account)
	if err != nil {
		logging.Error("Failed to sync rank", "puuid", accountPUUID, "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to sync rank")
		return
	}

	respondJSON(w, http.StatusOK, apiModels.SuccessResponse{
		Message: "Rank synced successfully",
	})
}
