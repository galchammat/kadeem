package riot

import (
	"fmt"
	"io"
	"kadeem/internal/logging"
)

func (r *RiotHandler) GetAccountID(region, gameName, tagLine string) (string, error) {
	// Placeholder implementation
	if gameName == "" || tagLine == "" {
		err := fmt.Errorf("gameName and tagLine cannot be empty")
		logging.Error(err.Error())
		return "", err
	}
	url := r.buildURL(region, fmt.Sprintf("/riot/account/v1/accounts/by-riot-id/%s/%s", gameName, tagLine))
	resp, err := r.makeRequest(url)
	if err != nil {
		logging.Error(err.Error())
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Error(err.Error())
		return "", err
	}
	// Simulate fetching account ID
	return string(body), nil
}
