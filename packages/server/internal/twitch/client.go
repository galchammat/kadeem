package twitch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/model"

	clientcredentials "golang.org/x/oauth2/clientcredentials"
)

type TwitchClient struct {
	ctx        context.Context
	httpClient *http.Client
	baseUrl    string
}

func NewTwitchClient(ctx context.Context) *TwitchClient {
	clientID := os.Getenv("TWITCH_CLIENT_ID")
	clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")

	conf := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     "https://id.twitch.tv/oauth2/token",
	}
	httpClient := conf.Client(ctx)

	// ensure transport exists and wrap it to inject Client-ID header required by Twitch
	if clientID != "" {
		base := httpClient.Transport
		if base == nil {
			base = http.DefaultTransport
		}
		httpClient.Transport = &clientIDTransport{base: base, clientID: clientID}
	}

	return &TwitchClient{
		ctx:        ctx,
		httpClient: httpClient,
		baseUrl:    "https://api.twitch.tv/helix",
	}
}

func (c *TwitchClient) buildURL(endpoint string) string {
	return fmt.Sprintf("%s%s", c.baseUrl, endpoint)
}

func (c *TwitchClient) makeRequest(endpoint string) (*model.TwitchResponse, int, error) {
	url := c.buildURL(endpoint)
	req, err := http.NewRequestWithContext(c.ctx, "GET", url, nil)
	if err != nil {
		logging.Error("Failed to create Twitch HTTP request", "endpoint", endpoint, "error", err)
		return nil, 0, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logging.Error("Twitch HTTP request failed", "url", url, "error", err)
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Error("Failed to read Twitch response body", "url", url, "error", err)
		return nil, 0, err
	}

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("HTTP request failed with status %d. body %s", resp.StatusCode, string(body))
		logging.Error("Twitch API returned non-200 status", "url", url, "statusCode", resp.StatusCode, "body", string(body))
		return nil, resp.StatusCode, err
	}

	var response model.TwitchResponse
	if err := json.Unmarshal(body, &response); err != nil {
		logging.Error("Failed to unmarshal Twitch response", "url", url, "error", err)
		return nil, resp.StatusCode, err
	}

	return &response, resp.StatusCode, nil
}

type clientIDTransport struct {
	base     http.RoundTripper
	clientID string
}

func (t *clientIDTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	r := req.Clone(req.Context())
	if t.clientID != "" {
		r.Header.Set("Client-ID", t.clientID)
	}
	return t.base.RoundTrip(r)
}
