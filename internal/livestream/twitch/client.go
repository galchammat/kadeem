package twitch

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/galchammat/kadeem/internal/logging"
	oauth2 "golang.org/x/oauth2/clientcredentials"
)

type TwitchClient struct {
	ctx        context.Context
	httpClient *http.Client
	baseUrl    string
}

func NewTwitchClient(ctx context.Context) *TwitchClient {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("TWITCH_CLIENT_ID"),
		ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
		TokenURL:     "https://id.twitch.tv/oauth2/token",
	}

	return &TwitchClient{
		ctx:        ctx,
		httpClient: conf.Client(ctx),
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
