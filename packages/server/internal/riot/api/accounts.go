package riot

import (
	"encoding/json"
	"fmt"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/model"
)

// FetchAccount fetches account data from the Riot API by gameName and tagLine.
func (c *Client) FetchAccount(region, gameName, tagLine string) (*model.LolAccount, error) {
	if gameName == "" || tagLine == "" || region == "" {
		return nil, fmt.Errorf("gameName, tagLine, and region cannot be empty")
	}

	url := c.buildURL(region, fmt.Sprintf("/riot/account/v1/accounts/by-riot-id/%s/%s", gameName, tagLine))
	body, statusCode, err := c.makeRequest(url)
	if err != nil {
		logging.Error("Failed to fetch account from Riot API", "gameName", gameName, "tagLine", tagLine, "region", region, "statusCode", statusCode, "error", err)
		if statusCode == 404 {
			return nil, fmt.Errorf("account not found on Riot servers")
		}
		return nil, fmt.Errorf("failed to fetch account from Riot servers: %v", err)
	}

	var account model.LolAccount
	if err := json.Unmarshal(body, &account); err != nil {
		logging.Error("Failed to unmarshal account JSON response", "gameName", gameName, "tagLine", tagLine, "error", err)
		return nil, err
	}
	account.Region = region
	return &account, nil
}

// FetchAccountByPUUID fetches account data from the Riot API by PUUID.
func (c *Client) FetchAccountByPUUID(region, puuid string) (*model.LolAccount, error) {
	url := c.buildURL(region, fmt.Sprintf("/riot/account/v1/accounts/by-puuid/%s", puuid))
	body, _, err := c.makeRequest(url)
	if err != nil {
		logging.Error("Failed to fetch account from Riot API", "puuid", puuid, "region", region, "error", err)
		return nil, err
	}

	var account model.LolAccount
	if err := json.Unmarshal(body, &account); err != nil {
		logging.Error("Failed to unmarshal account JSON response", "puuid", puuid, "error", err)
		return nil, err
	}
	account.Region = region
	return &account, nil
}
