package riot

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
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

type RiotHandler struct {
	ctx        context.Context
	httpClient *http.Client
	baseUrl    string
}

func NewRiotHandler(ctx context.Context) *RiotHandler {
	apiKey := os.Getenv("RIOT_API_KEY")

	// Create client with custom transport
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &riotTransport{
			apiKey: apiKey,
			base:   http.DefaultTransport,
		},
	}

	return &RiotHandler{
		ctx:        ctx,
		httpClient: client,
		baseUrl:    "https://{region}.api.riotgames.com",
	}
}

func (r *RiotHandler) buildURL(region, endpoint string) string {
	return fmt.Sprintf("https://%s.api.riotgames.com%s", region, endpoint)
}

func (r *RiotHandler) makeRequest(url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(r.ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// No need to set headers - transport handles it automatically!
	return r.httpClient.Do(req)
}
