package database

import (
	"strings"

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

// ListRiotAccounts lists accounts with optional filtering
func (db *DB) ListRiotAccounts(filter *models.LeagueOfLegendsAccount) ([]*models.LeagueOfLegendsAccount, error) {
	query := `SELECT puuid, tag_line, game_name, region FROM league_of_legends_accounts`
	var where []string
	var args []interface{}

	if filter != nil && filter.PUUID != "" {
		where = append(where, "puuid = ?")
		args = append(args, filter.PUUID)
	}
	if filter != nil && filter.TagLine != "" {
		where = append(where, "tag_line = ?")
		args = append(args, filter.TagLine)
	}
	if filter != nil && filter.GameName != "" {
		where = append(where, "game_name = ?")
		args = append(args, filter.GameName)
	}
	if filter != nil && filter.Region != "" {
		where = append(where, "region = ?")
		args = append(args, filter.Region)
	}

	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	rows, err := db.SQL.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*models.LeagueOfLegendsAccount
	for rows.Next() {
		account := &models.LeagueOfLegendsAccount{}
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

// DeleteRiotAccount deletes an account by PUUID
func (db *DB) DeleteRiotAccount(puuid string) error {
	query := `DELETE FROM league_of_legends_accounts WHERE puuid = ?`
	_, err := db.SQL.Exec(query, puuid)
	return err
}

// UpdateRiotAccount updates an existing account
func (db *DB) UpdateRiotAccount(account *models.LeagueOfLegendsAccount) error {
	query := `
		UPDATE league_of_legends_accounts 
		SET tag_line = ?, game_name = ?, region = ?
		WHERE puuid = ?`

	_, err := db.SQL.Exec(query, account.TagLine, account.GameName, account.Region, account.PUUID)
	return err
}
