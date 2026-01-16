package riot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/models"
)

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
		err := fmt.Errorf("replay download failed with status code %d", resp.StatusCode)
		logging.Error("Replay download failed with non-200 status", "matchID", matchID, "url", replayURL, "statusCode", resp.StatusCode)
		return err
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

func (c *RiotClient) SyncMatchReplay(matchID int64, url string) error {
	if err := downloadReplay(matchID, url); err != nil {
		return fmt.Errorf("error downloading replay: %v", err)
	}
	_, err := c.db.UpdateLolMatch(matchID, map[string]interface{}{"replay_synced": true})
	if err != nil {
		return err
	}
	return nil
}

func (c *RiotClient) FetchReplayURLs(puuid string, region string) ([]string, error) {
	endpoint := fmt.Sprintf("/lol/match/v5/matches/by-puuid/%s/replays", puuid)
	url := c.buildURL(region, endpoint)
	response, statusCode, err := c.makeRequest(url)
	if err != nil || statusCode != 200 {
		logging.Error("Failed to fetch replay URLs from Riot API", "puuid", puuid, "region", region, "statusCode", statusCode, "error", err)
		return nil, fmt.Errorf("error fetching replay URLs: %v Status Code: %d", err, statusCode)
	}
	var replays models.LolApiReplaysReponse
	if err := json.Unmarshal(response, &replays); err != nil {
		logging.Error("Failed to unmarshal replay URLs response", "puuid", puuid, "error", err)
		return nil, err
	}
	return replays.URLs, nil
}
