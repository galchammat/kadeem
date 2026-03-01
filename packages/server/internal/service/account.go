package service

import (
	"fmt"
	"reflect"

	"github.com/galchammat/kadeem/internal/database"
	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/model"
	riot "github.com/galchammat/kadeem/internal/riot/api"
)

type AccountService struct {
	db   *database.DB
	riot *riot.Client
}

func NewAccountService(db *database.DB, riot *riot.Client) *AccountService {
	return &AccountService{db: db, riot: riot}
}

// AddAccount fetches account from Riot API and saves it.
func (s *AccountService) AddAccount(region, gameName, tagLine string, streamerID int) error {
	account, err := s.riot.FetchAccount(region, gameName, tagLine)
	if err != nil {
		return err
	}
	account.StreamerID = streamerID
	return s.db.SaveRiotAccount(account)
}

// ReconcileAccount checks if account data has changed on Riot servers and updates DB.
func (s *AccountService) ReconcileAccount(account *model.LeagueOfLegendsAccount) error {
	fetched, err := s.riot.FetchAccountByPUUID(account.Region, account.PUUID)
	if err != nil {
		return err
	}
	fetched.Region = account.Region
	fetched.StreamerID = account.StreamerID

	if !reflect.DeepEqual(account, fetched) {
		logging.Info("Riot account data has changed, updating database", "puuid", account.PUUID)
		if err := s.db.SaveRiotAccount(fetched); err != nil {
			return err
		}
		*account = *fetched
	}
	return nil
}

// ListAccounts lists accounts with optional reconciliation.
func (s *AccountService) ListAccounts(filter *model.LeagueOfLegendsAccount) ([]model.LeagueOfLegendsAccount, error) {
	accounts, err := s.db.ListRiotAccounts(filter, 1000, 0)
	if err != nil {
		return nil, err
	}
	for i := range accounts {
		if err := s.ReconcileAccount(&accounts[i]); err != nil {
			return nil, err
		}
	}
	return accounts, nil
}

// UpdateAccount validates account against Riot API and updates DB.
func (s *AccountService) UpdateAccount(region, gameName, tagLine, puuid string) error {
	if gameName == "" || tagLine == "" || region == "" || puuid == "" {
		return fmt.Errorf("gameName, tagLine, region, and puuid cannot be empty")
	}

	validated, err := s.riot.FetchAccount(region, gameName, tagLine)
	if err != nil {
		return err
	}

	if validated.PUUID != puuid {
		return fmt.Errorf("the account %s#%s belongs to a different PUUID (%s), cannot update", gameName, tagLine, validated.PUUID)
	}

	_, err = s.db.UpdateRiotAccount(puuid, map[string]any{
		"game_name": validated.GameName,
		"tag_line":  validated.TagLine,
		"region":    region,
	})
	if err != nil {
		return err
	}

	logging.Info("Updated Riot account", "puuid", puuid)
	return nil
}

// DeleteAccount deletes a Riot account.
func (s *AccountService) DeleteAccount(puuid string) error {
	if puuid == "" {
		return fmt.Errorf("puuid cannot be empty")
	}
	if err := s.db.DeleteRiotAccount(puuid); err != nil {
		return err
	}
	logging.Info("Deleted Riot account", "puuid", puuid)
	return nil
}

// GetPlayerRankAtTime fetches the rank closest to a given timestamp.
func (s *AccountService) GetPlayerRankAtTime(puuid string, queueID int, timestamp int64) (*model.PlayerRank, error) {
	return s.db.GetRankAtTime(puuid, queueID, timestamp)
}
