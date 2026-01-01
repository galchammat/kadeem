package riot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/galchammat/kadeem/internal/models"
)

func replayStorageDir() string {
	if binDir := os.Getenv("BIN_DIR"); binDir != "" {
		return filepath.Join(binDir, "replays")
	}
	return filepath.Join("bin", "replays")
}

func downloadReplay(fileName string, replayURL string) error {
	if replayURL == "" {
		return fmt.Errorf("replay URL cannot be empty")
	}

	replaysDir := replayStorageDir()
	if err := os.MkdirAll(replaysDir, 0o755); err != nil {
		return err
	}

	filePath := filepath.Join(replaysDir, fileName+".rofl")
	resp, err := http.Get(replayURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("replay download failed with status code %d", resp.StatusCode)
	}

	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, resp.Body); err != nil {
		_ = outFile.Close()
		_ = os.Remove(filePath)
		return err
	}

	return nil
}

func (c *RiotClient) SyncMatchReplay(matchID string, url string) error {
	if err := downloadReplay(matchID, url); err != nil {
		return fmt.Errorf("error downloading replay: %v", err)
	}
	c.db.UpdateLolMatch(matchID, map[string]interface{}{"replay_synced": true})
	return nil
}

func (c *RiotClient) FetchReplayURLs(puuid string, region string) ([]string, error) {
	endpoint := fmt.Sprintf("/lol/match/v5/matches/by-puuid/%s/replays", puuid)
	url := c.buildURL(region, endpoint)
	response, statusCode, err := c.makeRequest(url)
	if err != nil || statusCode != 200 {
		return nil, fmt.Errorf("error fetching replay URLs: %v Status Code: %d", err, statusCode)
	}
	var replays models.LolApiReplaysReponse
	json.Unmarshal(response, &replays)
	return replays.URLs, nil
}
