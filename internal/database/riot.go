package database

import (
	"fmt"
	"strings"

	"github.com/galchammat/kadeem/internal/models"
)

// SaveRiotAccount saves a League of Legends account to the database
func (db *DB) SaveRiotAccount(account *models.LeagueOfLegendsAccount) error {
	query := `
        INSERT OR REPLACE INTO league_of_legends_accounts 
        (puuid, streamer_id, tag_line, game_name, region) 
        VALUES (?, ?, ?, ?, ?)`

	_, err := db.SQL.Exec(query, account.PUUID, account.StreamerID, account.TagLine, account.GameName, account.Region)
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
	if filter != nil && filter.StreamerID != 0 {
		where = append(where, "streamer_id = ?")
		args = append(args, filter.StreamerID)
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

func (db *DB) GetLolMatch(matchID string) (*models.LeagueOfLegendsMatchSummary, error) {
	query := `SELECT id, started_at, duration, replay_synced FROM league_of_legends_matches WHERE id = ?`

	var match models.LeagueOfLegendsMatchSummary
	err := db.SQL.QueryRow(query, matchID).Scan(&match.ID, &match.StartedAt, &match.Duration, &match.ReplaySynced)
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func (db *DB) InsertLolMatchSummary(summary *models.LeagueOfLegendsMatchSummary) error {
	if summary == nil || summary.StartedAt == nil || summary.Duration == nil {
		return fmt.Errorf("match summary is missing required fields")
	}

	replaySynced := false
	if summary.ReplaySynced != nil {
		replaySynced = *summary.ReplaySynced
	}

	query := `
		INSERT OR REPLACE INTO league_of_legends_matches
		(id, started_at, duration, replay_synced)
		VALUES (?, ?, ?, ?)`
	_, err := db.SQL.Exec(query, summary.ID, *summary.StartedAt, *summary.Duration, replaySynced)
	return err
}

func (db *DB) InsertLolMatchParticipantSummary(participant *models.LeagueOfLegendsMatchParticipantSummary) error {
	if participant == nil {
		return fmt.Errorf("participant summary cannot be nil")
	}

	query := `
		INSERT OR REPLACE INTO participants (
			game_id,
			champion_id,
			kills,
			deaths,
			assists,
			total_minions_killed,
			double_kills,
			triple_kills,
			quadra_kills,
			penta_kills,
			item0,
			item1,
			item2,
			item3,
			item4,
			item5,
			item6,
			summoner1_id,
			summoner2_id,
			lane,
			participant_id,
			puuid,
			riot_id_game_name,
			riot_id_tagline,
			total_damage_dealt_to_champions,
			total_damage_taken,
			win
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := db.SQL.Exec(
		query,
		participant.GameID,
		participant.ChampionID,
		participant.Kills,
		participant.Deaths,
		participant.Assists,
		participant.TotalMinionsKilled,
		participant.DoubleKills,
		participant.TripleKills,
		participant.QuadraKills,
		participant.PentaKills,
		participant.Item0,
		participant.Item1,
		participant.Item2,
		participant.Item3,
		participant.Item4,
		participant.Item5,
		participant.Item6,
		participant.Summoner1ID,
		participant.Summoner2ID,
		participant.Lane,
		participant.ParticipantID,
		participant.PUUID,
		participant.RiotIDGameName,
		participant.RiotIDTagline,
		participant.TotalDamageDealtToChampions,
		participant.TotalDamageTaken,
		participant.Win,
	)
	return err
}
