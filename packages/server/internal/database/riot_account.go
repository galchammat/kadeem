package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/model"
)

// SaveRiotAccount saves a League of Legends account to the database (shared pool)
func (db *DB) SaveRiotAccount(account *model.LeagueOfLegendsAccount) error {
	logging.Debug("updating account", "account", account)
	query := `
        INSERT INTO league_of_legends_accounts 
        (puuid, streamer_id, tag_line, game_name, region) 
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (puuid) DO UPDATE SET
			streamer_id = EXCLUDED.streamer_id,
			tag_line = EXCLUDED.tag_line,
			game_name = EXCLUDED.game_name,
			region = EXCLUDED.region`

	_, err := db.SQL.Exec(query, account.PUUID, account.StreamerID, account.TagLine, account.GameName, account.Region)
	if err != nil {
		logging.Error("Failed to save Riot account to database", "puuid", account.PUUID, "error", err)
		return err
	}
	return nil
}

// GetRiotAccount retrieves an account by PUUID
func (db *DB) GetRiotAccount(puuid string) (*model.LeagueOfLegendsAccount, error) {
	query := `SELECT puuid, tag_line, game_name, region, synced_at, streamer_id FROM league_of_legends_accounts WHERE puuid = $1`

	var account model.LeagueOfLegendsAccount
	err := db.SQL.QueryRow(query, puuid).Scan(&account.PUUID, &account.TagLine, &account.GameName, &account.Region, &account.SyncedAt, &account.StreamerID)
	if err != nil {
		logging.Error("Failed to get Riot account from database", "puuid", puuid, "error", err)
		return nil, err
	}

	return &account, nil
}

// GetRiotAccountByPUUID is an alias for GetRiotAccount (replaces former GetRiotAccountByID)
func (db *DB) GetRiotAccountByPUUID(puuid string) (*model.LeagueOfLegendsAccount, error) {
	return db.GetRiotAccount(puuid)
}

// FindRiotAccount finds an account by game name, tag line, and region
func (db *DB) FindRiotAccount(gameName, tagLine, region string) (*model.LeagueOfLegendsAccount, error) {
	query := `SELECT puuid, tag_line, game_name, region, synced_at, streamer_id FROM league_of_legends_accounts 
	          WHERE game_name = $1 AND tag_line = $2 AND region = $3`

	var account model.LeagueOfLegendsAccount
	err := db.SQL.QueryRow(query, gameName, tagLine, region).Scan(&account.PUUID, &account.TagLine, &account.GameName, &account.Region, &account.SyncedAt, &account.StreamerID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		logging.Error("Failed to find Riot account", "gameName", gameName, "tagLine", tagLine, "region", region, "error", err)
		return nil, err
	}

	return &account, nil
}

// FindOrCreateRiotAccount finds or creates an account (idempotent)
func (db *DB) FindOrCreateRiotAccount(gameName, tagLine, region string, streamerID int) (*model.LeagueOfLegendsAccount, error) {
	account, err := db.FindRiotAccount(gameName, tagLine, region)
	if err != nil {
		return nil, err
	}
	if account != nil {
		return account, nil
	}

	// Create new account
	newAccount := &model.LeagueOfLegendsAccount{
		GameName:   gameName,
		TagLine:    tagLine,
		Region:     region,
		StreamerID: streamerID,
	}
	err = db.SaveRiotAccount(newAccount)
	if err != nil {
		return nil, err
	}

	return newAccount, nil
}

// ListTrackedAccounts returns accounts a user is tracking with pagination
func (db *DB) ListTrackedAccounts(userID string, limit, offset int) ([]model.LeagueOfLegendsAccount, error) {
	query := `SELECT a.puuid, a.tag_line, a.game_name, a.region, a.synced_at, a.streamer_id 
	          FROM league_of_legends_accounts a
	          INNER JOIN user_tracked_accounts uta ON a.puuid = uta.account_puuid
	          WHERE uta.user_id = $1
	          ORDER BY uta.tracked_at DESC
	          LIMIT $2 OFFSET $3`

	rows, err := db.SQL.Query(query, userID, limit, offset)
	if err != nil {
		logging.Error("Failed to list tracked accounts", "userID", userID, "error", err)
		return nil, err
	}
	defer rows.Close()

	var accounts []model.LeagueOfLegendsAccount
	for rows.Next() {
		var account model.LeagueOfLegendsAccount
		if err := rows.Scan(&account.PUUID, &account.TagLine, &account.GameName, &account.Region, &account.SyncedAt, &account.StreamerID); err != nil {
			logging.Error("Failed to scan tracked account row", "error", err)
			return nil, err
		}
		accounts = append(accounts, account)
	}
	if err := rows.Err(); err != nil {
		logging.Error("Error iterating over tracked account rows", "error", err)
		return nil, err
	}
	return accounts, nil
}

// TrackAccount adds a tracking relationship (idempotent)
func (db *DB) TrackAccount(userID string, accountPUUID string) error {
	query := `INSERT INTO user_tracked_accounts (user_id, account_puuid) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := db.SQL.Exec(query, userID, accountPUUID)
	if err != nil {
		logging.Error("Failed to track account", "userID", userID, "accountPUUID", accountPUUID, "error", err)
	}
	return err
}

// UntrackAccount removes a tracking relationship
func (db *DB) UntrackAccount(userID string, accountPUUID string) error {
	query := `DELETE FROM user_tracked_accounts WHERE user_id = $1 AND account_puuid = $2`
	_, err := db.SQL.Exec(query, userID, accountPUUID)
	if err != nil {
		logging.Error("Failed to untrack account", "userID", userID, "accountPUUID", accountPUUID, "error", err)
	}
	return err
}

// IsTrackingAccount checks if user is tracking an account
func (db *DB) IsTrackingAccount(userID string, accountPUUID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM user_tracked_accounts WHERE user_id = $1 AND account_puuid = $2)`
	var exists bool
	err := db.SQL.QueryRow(query, userID, accountPUUID).Scan(&exists)
	if err != nil {
		logging.Error("Failed to check tracking status", "userID", userID, "accountPUUID", accountPUUID, "error", err)
		return false, err
	}
	return exists, nil
}

// GetTrackedAccountsForSync returns all accounts with at least one tracker (for background jobs)
func (db *DB) GetTrackedAccountsForSync() ([]model.LeagueOfLegendsAccount, error) {
	query := `SELECT DISTINCT a.puuid, a.tag_line, a.game_name, a.region, a.synced_at, a.streamer_id 
	          FROM league_of_legends_accounts a
	          INNER JOIN user_tracked_accounts uta ON a.puuid = uta.account_puuid`

	rows, err := db.SQL.Query(query)
	if err != nil {
		logging.Error("Failed to get tracked accounts for sync", "error", err)
		return nil, err
	}
	defer rows.Close()

	var accounts []model.LeagueOfLegendsAccount
	for rows.Next() {
		var account model.LeagueOfLegendsAccount
		if err := rows.Scan(&account.PUUID, &account.TagLine, &account.GameName, &account.Region, &account.SyncedAt, &account.StreamerID); err != nil {
			logging.Error("Failed to scan account row for sync", "error", err)
			return nil, err
		}
		accounts = append(accounts, account)
	}
	if err := rows.Err(); err != nil {
		logging.Error("Error iterating over accounts for sync", "error", err)
		return nil, err
	}
	return accounts, nil
}

// ListRiotAccounts lists accounts with optional filtering and pagination (for admin/internal use)
func (db *DB) ListRiotAccounts(filter *model.LeagueOfLegendsAccount, limit, offset int) ([]model.LeagueOfLegendsAccount, error) {
	query := `SELECT puuid, tag_line, game_name, region, synced_at, streamer_id FROM league_of_legends_accounts`
	var where []string
	var args []any
	argCounter := 1

	if filter != nil && filter.PUUID != "" {
		where = append(where, fmt.Sprintf("puuid = $%d", argCounter))
		args = append(args, filter.PUUID)
		argCounter++
	}
	if filter != nil && filter.TagLine != "" {
		where = append(where, fmt.Sprintf("tag_line = $%d", argCounter))
		args = append(args, filter.TagLine)
		argCounter++
	}
	if filter != nil && filter.GameName != "" {
		where = append(where, fmt.Sprintf("game_name = $%d", argCounter))
		args = append(args, filter.GameName)
		argCounter++
	}
	if filter != nil && filter.Region != "" {
		where = append(where, fmt.Sprintf("region = $%d", argCounter))
		args = append(args, filter.Region)
		argCounter++
	}
	if filter != nil && filter.StreamerID != 0 {
		where = append(where, fmt.Sprintf("streamer_id = $%d", argCounter))
		args = append(args, filter.StreamerID)
		argCounter++
	}

	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}
	query += fmt.Sprintf(" ORDER BY game_name LIMIT $%d OFFSET $%d", argCounter, argCounter+1)
	args = append(args, limit, offset)

	rows, err := db.SQL.Query(query, args...)
	if err != nil {
		logging.Error("Failed to list Riot accounts from database", "error", err)
		return nil, err
	}
	defer rows.Close()

	var accounts []model.LeagueOfLegendsAccount
	for rows.Next() {
		var account model.LeagueOfLegendsAccount
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

// DeleteRiotAccount deletes an account by PUUID (admin only)
func (db *DB) DeleteRiotAccount(puuid string) error {
	query := `DELETE FROM league_of_legends_accounts WHERE puuid = $1`
	_, err := db.SQL.Exec(query, puuid)
	if err != nil {
		logging.Error("Failed to delete Riot account from database", "puuid", puuid, "error", err)
	}
	return err
}

// allowedAccountColumns is the set of columns that can be updated via UpdateRiotAccount.
var allowedAccountColumns = map[string]bool{
	"tag_line":  true,
	"game_name": true,
	"region":    true,
	"synced_at": true,
}

func (db *DB) UpdateRiotAccount(PUUID string, updates map[string]any) (bool, error) {
	var setClauses []string
	var args []any
	argN := 1

	for column, value := range updates {
		if !allowedAccountColumns[column] {
			return false, fmt.Errorf("disallowed column: %s", column)
		}
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, argN))
		args = append(args, value)
		argN++
	}
	if len(setClauses) == 0 {
		return false, nil
	}
	args = append(args, PUUID)

	query := `UPDATE league_of_legends_accounts SET ` + strings.Join(setClauses, ", ") + fmt.Sprintf(` WHERE puuid = $%d`, argN)

	res, err := db.SQL.Exec(query, args...)
	if err != nil {
		logging.Error("Failed to update Riot account in database", "puuid", PUUID, "error", err)
		return false, err
	}
	n, _ := res.RowsAffected()
	return (n != 0), nil
}
