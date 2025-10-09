package database

import (
	"github.com/galchammat/kadeem/internal/models"
)

// SaveRiotAccount saves a League of Legends account to the database
func (db *DB) SaveRiotAccount(account *models.LeagueOfLegendsAccount) error {
	query := `
        INSERT OR REPLACE INTO league_of_legends_accounts 
        (puuid, tag_line, game_name, region) 
        VALUES (?, ?, ?, ?)`

	_, err := db.SQL.Exec(query, account.PUUID, account.TagLine, account.GameName, account.Region)
	return err
}

// GetRiotAccount retrieves an account by PUUID
func (db *DB) GetRiotAccount(puuid string) (*models.LeagueOfLegendsAccount, error) {
	query := `SELECT puuid, tag_line, game_name, region FROM league_of_legends_accounts WHERE puuid = ?`

	var account models.LeagueOfLegendsAccount
	err := db.SQL.QueryRow(query, puuid).Scan(&account.PUUID, &account.TagLine, &account.GameName, &account.Region)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// ListRiotAccountsByRegion gets all accounts for a specific region
func (db *DB) ListRiotAccountsByRegion(region string) ([]models.LeagueOfLegendsAccount, error) {
	query := `SELECT puuid, tag_line, game_name, region FROM league_of_legends_accounts WHERE region = ?`
	rows, err := db.SQL.Query(query, region)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.LeagueOfLegendsAccount
	for rows.Next() {
		var account models.LeagueOfLegendsAccount
		if err := rows.Scan(&account.PUUID, &account.TagLine, &account.GameName, &account.Region); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return accounts, nil
}

func (db *DB) DeleteRiotAccount(puuid string) error {
	_, err := db.SQL.Exec("DELETE FROM league_of_legends_accounts WHERE puuid = ?", puuid)
	return err
}
