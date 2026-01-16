package database

import (
	"strings"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/models"
)

// SaveRiotAccount saves a League of Legends account to the database
func (db *DB) SaveRiotAccount(account *models.LeagueOfLegendsAccount) error {
	logging.Debug("updating account", "account", account)
	query := `
        INSERT OR REPLACE INTO league_of_legends_accounts 
        (puuid, streamer_id, tag_line, game_name, region) 
        VALUES (?, ?, ?, ?, ?)`

	_, err := db.SQL.Exec(query, account.PUUID, account.StreamerID, account.TagLine, account.GameName, account.Region)
	if err != nil {
		logging.Error("Failed to save Riot account to database", "puuid", account.PUUID, "error", err)
	}
	return err
}

// GetRiotAccount retrieves an account by PUUID
func (db *DB) GetRiotAccount(puuid string) (*models.LeagueOfLegendsAccount, error) {
	query := `SELECT puuid, tag_line, game_name, region, synced_at, streamer_id FROM league_of_legends_accounts WHERE puuid = ?`

	var account models.LeagueOfLegendsAccount
	err := db.SQL.QueryRow(query, puuid).Scan(&account.PUUID, &account.TagLine, &account.GameName, &account.Region, &account.SyncedAt, &account.StreamerID)
	if err != nil {
		logging.Error("Failed to get Riot account from database", "puuid", puuid, "error", err)
		return nil, err
	}
	return &account, nil
}

// ListRiotAccounts lists accounts with optional filtering
func (db *DB) ListRiotAccounts(filter *models.LeagueOfLegendsAccount) ([]models.LeagueOfLegendsAccount, error) {
	query := `SELECT puuid, tag_line, game_name, region, synced_at, streamer_id FROM league_of_legends_accounts`
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
	if filter != nil && filter.StreamerID != 0 {
		where = append(where, "streamer_id = ?")
		args = append(args, filter.StreamerID)
	}

	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	rows, err := db.SQL.Query(query, args...)
	if err != nil {
		logging.Error("Failed to list Riot accounts from database", "error", err)
		return nil, err
	}
	defer rows.Close()

	var accounts []models.LeagueOfLegendsAccount
	for rows.Next() {
		var account models.LeagueOfLegendsAccount
		if err := rows.Scan(&account.PUUID, &account.TagLine, &account.GameName, &account.Region, &account.SyncedAt, &account.StreamerID); err != nil {
			logging.Error("Failed to scan Riot account row", "error", err)
			return nil, err
		}

		accounts = append(accounts, account)
	}
	if err := rows.Err(); err != nil {
		logging.Error("Error iterating over Riot account rows", "error", err)
		return nil, err
	}
	return accounts, nil
}

// DeleteRiotAccount deletes an account by PUUID
func (db *DB) DeleteRiotAccount(puuid string) error {
	query := `DELETE FROM league_of_legends_accounts WHERE puuid = ?`
	_, err := db.SQL.Exec(query, puuid)
	if err != nil {
		logging.Error("Failed to delete Riot account from database", "puuid", puuid, "error", err)
	}
	return err
}

func (db *DB) UpdateRiotAccount(PUUID string, updates map[string]interface{}) (bool, error) {
	var setClauses []string
	var args []interface{}

	for column, value := range updates {
		setClauses = append(setClauses, column+" = ?")
		args = append(args, value)
	}
	args = append(args, PUUID)

	query := `UPDATE league_of_legends_accounts SET ` + strings.Join(setClauses, ", ") + ` WHERE puuid = ?`

	res, err := db.SQL.Exec(query, args...)
	if err != nil {
		logging.Error("Failed to update Riot account in database", "puuid", PUUID, "error", err)
		return false, err
	}
	n, _ := res.RowsAffected()
	return (n != 0), nil
}
