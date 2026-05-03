package riot

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"

	"github.com/galchammat/kadeem/internal/model"
	"github.com/galchammat/kadeem/internal/syncer"
)

var _ syncer.ArtifactHandler = (*ReplayHandler)(nil)
var _ syncer.ArtifactStore = (*ReplayStore)(nil)

type ReplayWriter interface {
	WriteReplay(ctx context.Context, key string, body io.Reader) error
}

type ReplayHandler struct {
	client   *Client
	writer   ReplayWriter
	accounts []model.LolAccount
}

func NewReplayHandler(client *Client, writer ReplayWriter, accounts []model.LolAccount) (*ReplayHandler, error) {
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

type ReplayStore struct {
	db *sql.DB
}

func NewReplayStore(db *sql.DB) (*ReplayStore, error) {
	if db == nil {
		return nil, fmt.Errorf("replay db is nil")
	}

	return &ReplayStore{db: db}, nil
}

func (s *ReplayStore) ClaimPending(ctx context.Context, limit int) ([]syncer.Artifact, error) {
	if limit <= 0 {
		return nil, nil
	}

	rows, err := s.db.QueryContext(ctx, `
		WITH pending AS (
			SELECT id
			FROM lol_matches
			WHERE replay_s3_key IS NULL
			ORDER BY id DESC
			LIMIT $1
			FOR UPDATE SKIP LOCKED
		)
		UPDATE lol_matches
		SET replay_sync_attempted_at = NOW(),
			replay_sync_error = NULL
		FROM pending
		WHERE lol_matches.id = pending.id
		RETURNING lol_matches.id
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("claim pending replays: %w", err)
	}
	defer rows.Close()

	artifacts := make([]syncer.Artifact, 0, limit)
	for rows.Next() {
		var matchID int64
		if err := rows.Scan(&matchID); err != nil {
			return nil, fmt.Errorf("scan replay artifact: %w", err)
		}

		id := strconv.FormatInt(matchID, 10)
		artifacts = append(artifacts, syncer.Artifact{
			ID:         id,
			ExternalID: id,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate replay artifacts: %w", err)
	}

	return artifacts, nil
}

func (s *ReplayStore) MarkDone(ctx context.Context, id string, s3Key string) error {
	matchID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("parse replay artifact id %q: %w", id, err)
	}

	result, err := s.db.ExecContext(ctx, `
		UPDATE lol_matches
		SET replay_s3_key = $2,
			replay_sync_error = NULL
		WHERE id = $1
	`, matchID, s3Key)
	if err != nil {
		return fmt.Errorf("mark replay done %q: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read replay mark done result %q: %w", id, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("replay artifact %q not found", id)
	}

	return nil
}

func (s *ReplayStore) MarkFailed(ctx context.Context, id string, err error) error {
	matchID, parseErr := strconv.ParseInt(id, 10, 64)
	if parseErr != nil {
		return fmt.Errorf("parse replay artifact id %q: %w", id, parseErr)
	}

	_, updateErr := s.db.ExecContext(ctx, `
		UPDATE lol_matches
		SET replay_sync_error = $2,
			replay_sync_attempted_at = NOW()
		WHERE id = $1
	`, matchID, err.Error())
	if updateErr != nil {
		return fmt.Errorf("mark replay failed %q: %w", id, updateErr)
	}

	return nil
}

var replayURLPattern = regexp.MustCompile(`([A-Z0-9]+_([0-9]+))\.replay`)

func replayMatchID(url string) (string, error) {
	matches := replayURLPattern.FindStringSubmatch(url)
	if len(matches) < 3 {
		return "", fmt.Errorf("replay match id not found in url %q", url)
	}

	return matches[2], nil
}
