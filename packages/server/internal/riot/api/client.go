package riot

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/galchammat/kadeem/internal/logging"
)

// riotTransport adds API key to all requests.
type riotTransport struct {
	apiKey string
	base   http.RoundTripper
}

func (t *riotTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.Header.Set("X-Riot-Token", t.apiKey)
	req.Header.Set("Content-Type", "application/json")
	if t.base == nil {
		t.base = http.DefaultTransport
	}
	return t.base.RoundTrip(req)
}

// Client is a pure HTTP client for the Riot Games API.
// It has no database dependency.
type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	apiKey := os.Getenv("RIOT_API_KEY")
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &riotTransport{
				apiKey: apiKey,
				base:   http.DefaultTransport,
			},
		},
	}
}

func (c *Client) buildURL(region, endpoint string) string {
	generalRegion, err := GetAPIRegion(region)
	if err != nil {
		logging.Error("Failed to generalize API region", "region", region)
	}
	return fmt.Sprintf("https://%s.api.riotgames.com%s", generalRegion, endpoint)
}

func (c *Client) makeRequest(url string) ([]byte, int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 400, err
	}

	resp, err := c.httpClient.Do(req)
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

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("HTTP request failed. url %s. status %d. body %s", url, resp.StatusCode, string(body))
		logging.Error(err.Error())
		return nil, resp.StatusCode, err
	}

	return body, resp.StatusCode, nil
}
