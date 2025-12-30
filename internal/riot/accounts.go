package riot

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/models"
)

func (r *RiotClient) AddAccount(region, gameName, tagLine string, streamerID int) error {
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
	body, statusCode, err := r.makeRequest(url)
	if err != nil {
		logging.Error(err.Error())
		if statusCode == 404 {
			err := fmt.Errorf("account not found on Riot servers")
			return err
		} else {
			return fmt.Errorf("failed to fetch account from Riot servers: %v", err)
		}
	}
	if err := json.Unmarshal(body, &account); err != nil {
		logging.Error("Failed to unmarshal JSON response", "error", err)
		return err
	}
	account.Region = region
	account.StreamerID = streamerID
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
	body, _, err := r.makeRequest(url)
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

func (r *RiotClient) DeleteAccount(puuid string) error {
	if puuid == "" {
		err := fmt.Errorf("puuid cannot be empty")
		logging.Error(err.Error())
		return err
	}
	if err := r.db.DeleteRiotAccount(puuid); err != nil {
		logging.Error("Failed to delete Riot account: %v", err)
		return err
	}
	logging.Info("Deleted Riot account", "puuid", puuid)
	return nil
}

func (r *RiotClient) UpdateAccount(region, gameName, tagLine, puuid string) error {
	if gameName == "" || tagLine == "" || region == "" || puuid == "" {
		err := fmt.Errorf("gameName, tagLine, region, and puuid cannot be empty")
		logging.Error(err.Error())
		return err
	}

	// Validate the account exists on Riot's servers
	apiRegion, err := GetAPIRegion(region)
	if err != nil {
		logging.Error(err.Error())
		return err
	}
	url := r.buildURL(apiRegion, fmt.Sprintf("/riot/account/v1/accounts/by-riot-id/%s/%s", gameName, tagLine))
	body, statusCode, err := r.makeRequest(url)
	if err != nil {
		if statusCode == 404 {
			return fmt.Errorf("account not found on Riot servers")
		}
		return fmt.Errorf("failed to fetch account from Riot servers: %v", err)
	}

	var validatedAccount models.LeagueOfLegendsAccount
	if err := json.Unmarshal(body, &validatedAccount); err != nil {
		logging.Error("Failed to unmarshal JSON response", "error", err)
		return err
	}

	// Check if the validated account PUUID matches
	if validatedAccount.PUUID != puuid {
		return fmt.Errorf("the account %s#%s belongs to a different PUUID (%s), cannot update", gameName, tagLine, validatedAccount.PUUID)
	}

	validatedAccount.Region = region
	if err := r.db.UpdateRiotAccount(&validatedAccount); err != nil {
		logging.Error("Failed to update Riot account: %v", err)
		return err
	}

	logging.Info("Updated Riot account", "puuid", puuid)
	return nil
}
