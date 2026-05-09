package riot

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"

	riotapi "github.com/galchammat/kadeem/internal/riot/api"
	"github.com/galchammat/kadeem/internal/riot/models"
	"github.com/galchammat/kadeem/internal/syncer"
)

var _ syncer.ArtifactHandler = (*ReplayHandler)(nil)

type ReplayWriter interface {
	WriteReplay(ctx context.Context, key string, body io.Reader) error
}

type ReplayHandler struct {
	client   *riotapi.Client
	writer   ReplayWriter
	accounts []models.Account
}

func NewReplayHandler(client *riotapi.Client, writer ReplayWriter, accounts []models.Account) (*ReplayHandler, error) {
	if client == nil {
		return nil, fmt.Errorf("riot client is nil")
	}
	if writer == nil {
		return nil, fmt.Errorf("replay writer is nil")
	}

	return &ReplayHandler{
		client:   client,
		writer:   writer,
		accounts: accounts,
	}, nil
}

func (r *ReplayHandler) Process(ctx context.Context, artifact syncer.Artifact) (string, error) {
	if artifact.ExternalID == "" {
		return "", fmt.Errorf("artifact external id is empty")
	}

	replayURL, err := r.replayURL(ctx, artifact.ExternalID)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, replayURL, nil)
	if err != nil {
		return "", fmt.Errorf("create replay request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("download replay %q: %w", artifact.ExternalID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download replay %q: status %d", artifact.ExternalID, resp.StatusCode)
	}

	key := fmt.Sprintf("lol/replays/%s.rofl", artifact.ExternalID)
	if err := r.writer.WriteReplay(ctx, key, resp.Body); err != nil {
		return "", fmt.Errorf("store replay %q: %w", artifact.ExternalID, err)
	}

	return key, nil
}

func (r *ReplayHandler) replayURL(ctx context.Context, matchID string) (string, error) {
	for _, account := range r.accounts {
		if err := ctx.Err(); err != nil {
			return "", err
		}

		urls, err := r.client.FetchReplayURLs(account.PUUID, account.Region)
		if err != nil {
			return "", fmt.Errorf("fetch replay urls for puuid %q: %w", account.PUUID, err)
		}

		for _, url := range urls {
			urlMatchID, err := replayMatchID(url)
			if err != nil {
				continue
			}
			if urlMatchID == matchID {
				return url, nil
			}
		}
	}

	return "", fmt.Errorf("replay url not found for match %q", matchID)
}

var replayURLPattern = regexp.MustCompile(`([A-Z0-9]+_([0-9]+))\.replay`)

func replayMatchID(url string) (string, error) {
	matches := replayURLPattern.FindStringSubmatch(url)
	if len(matches) < 3 {
		return "", fmt.Errorf("replay match id not found in url %q", url)
	}

	return matches[2], nil
}
