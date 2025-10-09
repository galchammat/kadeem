package riot

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/models"
)

func (r *RiotClient) AddAccount(region, gameName, tagLine string) error {
	var account models.LeagueOfLegendsAccount
	if gameName == "" || tagLine == "" || region == "" {
		err := fmt.Errorf("gameName, tagLine, and region cannot be empty")
		logging.Error(err.Error())
		return err
	}
	apiRegion, err := GetAPIRegion(region)
	if err != nil {
		logging.Error(err.Error())
		return err
	}
	url := r.buildURL(apiRegion, fmt.Sprintf("/riot/account/v1/accounts/by-riot-id/%s/%s", gameName, tagLine))
	body, err := r.makeRequest(url)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &account); err != nil {
		logging.Error("Failed to unmarshal JSON response", "error", err)
		return err
	}
	account.Region = region
	if err := r.db.SaveRiotAccount(&account); err != nil {
		logging.Error("Failed to save Riot account: %v", err)
		return err
	}
	return nil
}

func (r *RiotClient) reconcileAccount(account *models.LeagueOfLegendsAccount) error {
	apiRegion, err := GetAPIRegion(account.Region)
	if err != nil {
		logging.Error(err.Error())
		return err
	}
	url := r.buildURL(apiRegion, fmt.Sprintf("/riot/account/v1/accounts/by-puuid/%s", account.PUUID))
	body, err := r.makeRequest(url)
	if err != nil {
		return err
	}
	var fetchedAccount models.LeagueOfLegendsAccount
	if err := json.Unmarshal(body, &fetchedAccount); err != nil {
		logging.Error("Failed to unmarshal JSON response", "error", err)
		return err
	}
	fetchedAccount.Region = account.Region

	if !reflect.DeepEqual(account, &fetchedAccount) {
		logging.Info("Riot account data has changed, updating database", "puuid", account.PUUID)
		if err := r.db.SaveRiotAccount(&fetchedAccount); err != nil {
			logging.Error("Failed to update Riot account: %v", err)
			return err
		}
		*account = fetchedAccount
	} else {
		logging.Debug("Riot account data is up-to-date", "name", account.GameName, "tag", account.TagLine)
	}
	return nil
}

func (r *RiotClient) ListAccounts(filter *models.LeagueOfLegendsAccount) ([]models.LeagueOfLegendsAccount, error) {
	accounts, err := r.db.ListRiotAccounts(filter)
	if err != nil {
		logging.Error("Failed to list Riot accounts while fetching puuids from database: %v", err)
		return nil, err
	}

	for _, account := range accounts {
		err := r.reconcileAccount(account)
		if err != nil {
			logging.Error("Failed to reconcile Riot account: %v", err)
			return nil, err
		}
	}

	result := make([]models.LeagueOfLegendsAccount, len(accounts))
	for i, acc := range accounts {
		result[i] = *acc
	}
	return result, nil
}
