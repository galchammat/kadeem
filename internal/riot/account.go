package riot

import (
	"encoding/json"
	"fmt"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/models"
)

func (r *RiotClient) GetAccount(region, gameName, tagLine string) (models.LeagueOfLegendsAccount, error) {
	var account models.LeagueOfLegendsAccount
	if gameName == "" || tagLine == "" || region == "" {
		err := fmt.Errorf("gameName, tagLine, and region cannot be empty")
		logging.Error(err.Error())
		return account, err
	}
	apiRegion, err := GetAPIRegion(region)
	if err != nil {
		logging.Error(err.Error())
		return account, err
	}
	url := r.buildURL(apiRegion, fmt.Sprintf("/riot/account/v1/accounts/by-riot-id/%s/%s", gameName, tagLine))
	body, err := r.makeRequest(url)
	if err != nil {
		return account, err
	}
	if err := json.Unmarshal(body, &account); err != nil {
		logging.Error("Failed to unmarshal JSON response", "error", err)
		return account, err
	}
	return account, nil
}
