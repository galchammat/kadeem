package service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/galchammat/kadeem/internal/constants"
	"github.com/galchammat/kadeem/internal/database"
	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/model"
	riot "github.com/galchammat/kadeem/internal/riot/api"
)

type MatchService struct {
	db   *database.DB
	riot *riot.Client
}

func NewMatchService(db *database.DB, riot *riot.Client) *MatchService {
	return &MatchService{db: db, riot: riot}
}

// SyncMatches syncs replays and match summaries for an account.
func (s *MatchService) SyncMatches(account model.LeagueOfLegendsAccount) error {
	logging.Debug("Syncing matches for account", "ID", account.PUUID)

	replayURLs, err := s.riot.FetchReplayURLs(account.PUUID, account.Region)
	if err != nil {
		logging.Error("Failed to fetch replay URLs for account", "puuid", account.PUUID, "region", account.Region, "error", err)
		return err
	}

	for _, url := range replayURLs {
		matchID, fullMatchID, err := extractMatchID(url)
		if err != nil {
			return fmt.Errorf("failed to parse matchID from replay URL: %s", url)
		}

		existingMatches, err := s.db.ListLolMatches(&model.LolMatchFilter{MatchID: &matchID}, 1, 0)
		if err != nil {
			return fmt.Errorf("error while checking for an existing match. MatchID: %d. Error: %w", matchID, err)
		}
		var existingMatch *model.LeagueOfLegendsMatch
		if len(existingMatches) != 0 {
			existingMatch = &existingMatches[0]
		}

		// Fetch the match summary if record does not exist or has no start timestamp
		if existingMatch == nil || existingMatch.Summary.StartedAt == nil {
			logging.Debug("Fetching match summary", "MatchID", matchID, "FullMatchID", fullMatchID)
			if err := s.SyncMatchSummary(matchID, fullMatchID, account.Region); err != nil {
				logging.Warn("Skipping match summary sync due to error", "MatchID", matchID)
			}
		}

		// Download the replay if record does not exist or has no replay
		if existingMatch == nil || existingMatch.ReplayURL == nil {
			logging.Debug("Downloading replay", "MatchID", matchID, "URL", url)
			if err := s.SyncMatchReplay(matchID, url); err != nil {
				logging.Warn("Skipping replay download due to error", "MatchID", matchID)
			}
		}
	}

	// Update sync timestamp
	_, err = s.db.UpdateRiotAccount(account.PUUID, map[string]any{"synced_at": time.Now().Unix()})
	return err
}

// SyncMatchSummary fetches match detail from Riot API and stores it.
func (s *MatchService) SyncMatchSummary(matchID int64, fullMatchID, region string) error {
	if matchID == 0 {
		return fmt.Errorf("matchID cannot be zero")
	}

	response, err := s.riot.FetchMatchDetail(fullMatchID, region)
	if err != nil {
		return err
	}

	summary := model.LeagueOfLegendsMatchSummary{
		ID:        response.Info.ID,
		StartedAt: &response.Info.StartedAt,
		Duration:  &response.Info.Duration,
	}

	if err := s.db.InsertLolMatchWithParticipants(&summary, response.Info.Participants); err != nil {
		logging.Error(
			"Failed to insert match with participants (transaction rolled back)",
			"matchID", matchID,
			"fullMatchID", fullMatchID,
			"participantCount", len(response.Info.Participants),
			"error", err,
		)
		return err
	}

	logging.Debug(
		"Successfully synced match summary with participants",
		"matchID", matchID,
		"participantCount", len(response.Info.Participants),
	)
	return nil
}

// SyncMatchReplay downloads a replay file and marks it synced in DB.
func (s *MatchService) SyncMatchReplay(matchID int64, replayURL string) error {
	if err := downloadReplay(matchID, replayURL); err != nil {
		return fmt.Errorf("error downloading replay: %v", err)
	}
	_, err := s.db.UpdateLolMatch(matchID, map[string]any{"replay_synced": true})
	return err
}

// ListMatches lists matches with auto-sync if stale.
func (s *MatchService) ListMatches(filter *model.LolMatchFilter, account *model.LeagueOfLegendsAccount, limit, offset int) ([]model.LeagueOfLegendsMatch, error) {
	if account != nil &&
		(account.SyncedAt == nil || time.Since(time.Unix(*account.SyncedAt, 0)) > constants.SyncRefreshInMinutes*time.Minute) {
		if err := s.SyncMatches(*account); err != nil {
			logging.Warn("Failed to sync matches for account, returning cached data", "PUUID", account.PUUID)
		}
	}
	return s.db.ListLolMatches(filter, limit, offset)
}

// FetchMatchIDs fetches match IDs from the Riot API.
func (s *MatchService) FetchMatchIDs(puuid, region string, startTime *int64) ([]string, error) {
	return s.riot.FetchMatchIDs(puuid, region, startTime)
}

// FetchReplayURLs fetches replay URLs from the Riot API.
func (s *MatchService) FetchReplayURLs(puuid, region string) ([]string, error) {
	return s.riot.FetchReplayURLs(puuid, region)
}

// --- helpers ---

func extractMatchID(url string) (int64, string, error) {
	re := regexp.MustCompile(`([A-Z0-9]+_\d+)\.replay`)
	matches := re.FindStringSubmatch(url)

	if len(matches) < 2 {
		return 0, "", fmt.Errorf("no match ID found in URL: %s", url)
	}

	fullMatchID := matches[1]
	parts := regexp.MustCompile(`_(\d+)$`).FindStringSubmatch(fullMatchID)
	if len(parts) < 2 {
		return 0, "", fmt.Errorf("invalid match ID format: %s", fullMatchID)
	}

	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, "", fmt.Errorf("invalid match ID in URL: %s", url)
	}

	return id, fullMatchID, nil
}

func replayStorageDir() string {
	if binDir := os.Getenv("BIN_DIR"); binDir != "" {
		return filepath.Join(binDir, "replays")
	}
	return filepath.Join("bin", "replays")
}

func replayExists(filePath string) bool {
	info, err := os.Stat(filePath)
	return err == nil && info.Size() > 1024*1024 // > 1MB
}

func downloadReplay(matchID int64, replayURL string) error {
	if replayURL == "" {
		return fmt.Errorf("replay URL cannot be empty")
	}

	replaysDir := replayStorageDir()
	if err := os.MkdirAll(replaysDir, 0o755); err != nil {
		logging.Error("Failed to create replays directory", "path", replaysDir, "error", err)
		return err
	}

	filePath := filepath.Join(replaysDir, fmt.Sprintf("%d.rofl", matchID))
	if replayExists(filePath) {
		return nil
	}

	resp, err := http.Get(replayURL)
	if err != nil {
		logging.Error("Failed to download replay from URL", "matchID", matchID, "url", replayURL, "error", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("replay download failed with status code %d", resp.StatusCode)
	}

	outFile, err := os.Create(filePath)
	if err != nil {
		logging.Error("Failed to create replay file", "matchID", matchID, "path", filePath, "error", err)
		return err
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, resp.Body); err != nil {
		logging.Error("Failed to write replay data to file", "matchID", matchID, "path", filePath, "error", err)
		_ = outFile.Close()
		_ = os.Remove(filePath)
		return err
	}

	return nil
}
