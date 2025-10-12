package riot

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/galchammat/kadeem/internal/database"
	"github.com/galchammat/kadeem/internal/logging"
)

// Custom RoundTripper that adds API key to all requests
type riotTransport struct {
	apiKey string
	base   http.RoundTripper
}

func (t *riotTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original
	req = req.Clone(req.Context())

	// Add API key header to every request
	req.Header.Set("X-Riot-Token", t.apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Use the base transport (or default if nil)
	if t.base == nil {
		t.base = http.DefaultTransport
	}

	return t.base.RoundTrip(req)
}

type RiotClient struct {
	ctx        context.Context
	httpClient *http.Client
	baseUrl    string
	db         *database.DB
}

func NewRiotClient(ctx context.Context, db *database.DB) *RiotClient {
	apiKey := os.Getenv("RIOT_API_KEY")

	// Create client with custom transport
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &riotTransport{
			apiKey: apiKey,
			base:   http.DefaultTransport,
		},
	}

	return &RiotClient{
		ctx:        ctx,
		httpClient: client,
		baseUrl:    "https://{region}.api.riotgames.com",
		db:         db,
	}
}

func (r *RiotClient) buildURL(region, endpoint string) string {
	return fmt.Sprintf("https://%s.api.riotgames.com%s", region, endpoint)
}

func (r *RiotClient) makeRequest(url string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(r.ctx, "GET", url, nil)
	if err != nil {
		return nil, 400, err
	}

	resp, err := r.httpClient.Do(req)
	// HTTP Error
	if err != nil {
		logging.Error(err.Error())
		return nil, 500, err
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		logging.Error(err.Error())
		return nil, 500, err
	}

	// Non-200 Status Code
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("HTTP request failed with status %d. body %s", resp.StatusCode, string(body))
		logging.Error(err.Error())
		return nil, resp.StatusCode, err
	}

	return body, resp.StatusCode, nil
}
