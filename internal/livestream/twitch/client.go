package twitch

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/galchammat/kadeem/internal/logging"
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

func (c *TwitchClient) makeRequest(endpoint string) ([]byte, int, error) {
	url := c.buildURL(endpoint)
	req, err := http.NewRequestWithContext(c.ctx, "GET", url, nil)
	if err != nil {
		return nil, 400, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logging.Error(err.Error())
		return nil, 500, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Error(err.Error())
		return nil, 500, err
	}

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("HTTP request failed with status %d. body %s", resp.StatusCode, string(body))
		logging.Error(err.Error())
		return nil, resp.StatusCode, err
	}

	return body, resp.StatusCode, nil
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
