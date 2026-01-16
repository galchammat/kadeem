package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/models"
)

// SQLExecutor interface allows functions to accept either *sql.DB or *sql.Tx
type SQLExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

func (db *DB) InsertLolMatchSummary(summary *models.LeagueOfLegendsMatchSummary) error {
	return db.insertLolMatchSummaryExec(db.SQL, summary)
}

func (db *DB) insertLolMatchSummaryExec(exec SQLExecutor, summary *models.LeagueOfLegendsMatchSummary) error {
	if summary == nil || summary.StartedAt == nil || summary.Duration == nil {
		return fmt.Errorf("match summary is missing required fields")
	}

	replaySynced := false
	if summary.ReplaySynced != nil {
		replaySynced = *summary.ReplaySynced
	}

	query := `
		INSERT OR REPLACE INTO lol_matches
		(id, started_at, duration, replay_synced)
		VALUES (?, ?, ?, ?)`
	_, err := exec.Exec(query, summary.ID, *summary.StartedAt, *summary.Duration, replaySynced)
	if err != nil {
		logging.Error("Failed to insert match summary into database", "matchID", summary.ID, "error", err)
	}
	return err
}

func (db *DB) UpdateLolMatch(matchID int64, updates map[string]interface{}) (bool, error) {
	var setClauses []string
	var args []interface{}

	for column, value := range updates {
		setClauses = append(setClauses, column+" = ?")
		args = append(args, value)
	}
	args = append(args, matchID)

	query := `UPDATE lol_matches SET ` + strings.Join(setClauses, ", ") + ` WHERE id = ?`

	res, err := db.SQL.Exec(query, args...)
	if err != nil {
		logging.Error("Failed to update match in database", "matchID", matchID, "error", err)
		return false, err
	}
	n, _ := res.RowsAffected()
	return (n != 0), nil
}

// InsertLolMatchWithParticipants atomically inserts a match summary and all its participants
// in a single transaction. If any insert fails, the entire transaction is rolled back.
func (db *DB) InsertLolMatchWithParticipants(
	summary *models.LeagueOfLegendsMatchSummary,
	participants []models.LeagueOfLegendsMatchParticipantSummary,
) error {
	tx, err := db.SQL.Begin()
	if err != nil {
		logging.Error("Failed to begin transaction for match insert", "matchID", summary.ID, "error", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Safe to call even after Commit

	// Insert match summary
	if err := db.insertLolMatchSummaryExec(tx, summary); err != nil {
		logging.Error("Failed to insert match summary, rolling back transaction", "matchID", summary.ID, "error", err)
		return err
	}

	// Insert all participants
	for i, participant := range participants {
		if err := db.insertLolMatchParticipantSummaryExec(tx, &participant); err != nil {
			logging.Error(
				"Failed to insert participant, rolling back entire match transaction",
				"matchID", summary.ID,
				"participantIndex", i,
				"participantID", participant.ParticipantID,
				"error", err,
			)
			return err
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		logging.Error("Failed to commit match transaction", "matchID", summary.ID, "error", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logging.Debug("Successfully inserted match with participants", "matchID", summary.ID, "participantCount", len(participants))
	return nil
}

func (db *DB) InsertLolMatchParticipantSummary(participant *models.LeagueOfLegendsMatchParticipantSummary) error {
	return db.insertLolMatchParticipantSummaryExec(db.SQL, participant)
}

func (db *DB) insertLolMatchParticipantSummaryExec(exec SQLExecutor, participant *models.LeagueOfLegendsMatchParticipantSummary) error {
	if participant == nil {
		return fmt.Errorf("participant summary cannot be nil")
	}

	query := `
		INSERT OR REPLACE INTO participants (
			match_id,
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

	_, err := exec.Exec(
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
	if err != nil {
		logging.Error("Failed to insert match participant into database", "matchID", participant.GameID, "participantID", participant.ParticipantID, "error", err)
	}
	return err
}

func (db *DB) ListLolMatches(filter *models.LolMatchFilter, limit *int, offset *int) ([]models.LeagueOfLegendsMatch, error) {
	// Default and max limit to 100
	if limit == nil || *limit <= 0 || *limit > 100 {
		*limit = 100
	}

	// Step 1: Build query to get matching match IDs
	matchIDQuery := `
		SELECT DISTINCT m.id
		FROM lol_matches m
		LEFT JOIN participants p ON m.id = p.match_id
	`

	// Build WHERE clauses using BuildQueryArgs
	whereClauses, args, err := db.BuildQueryArgs(filter)
	if err != nil {
		logging.Error("Failed to build query args for ListLolMatches", "error", err)
		return nil, fmt.Errorf("failed to build query args: %w", err)
	}

	// Add WHERE clause
	if len(whereClauses) > 0 {
		matchIDQuery += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Add ORDER BY, LIMIT, OFFSET
	matchIDQuery += " ORDER BY m.started_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	// Step 2: Execute first query to get match IDs
	rows, err := db.SQL.Query(matchIDQuery, args...)
	if err != nil {
		logging.Error("Failed to query match IDs from database", "error", err)
		return nil, err
	}
	defer rows.Close()

	var matchIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			logging.Error("Failed to scan match ID row", "error", err)
			return nil, err
		}
		matchIDs = append(matchIDs, id)
	}
	if err := rows.Err(); err != nil {
		logging.Error("Error iterating over match ID rows", "error", err)
		return nil, err
	}

	// If no matches found, return empty slice
	if len(matchIDs) == 0 {
		return []models.LeagueOfLegendsMatch{}, nil
	}

	// Step 3: Build second query to get full match data with all participants
	placeholders := make([]string, len(matchIDs))
	fullArgs := make([]interface{}, len(matchIDs))
	for i, id := range matchIDs {
		placeholders[i] = "?"
		fullArgs[i] = id
	}

	fullQuery := fmt.Sprintf(`
		SELECT 
			m.id, m.started_at, m.duration, m.replay_synced,
			p.match_id, p.champion_id, p.kills, p.deaths, p.assists,
			p.total_minions_killed, p.double_kills, p.triple_kills,
			p.quadra_kills, p.penta_kills, p.item0, p.item1, p.item2,
			p.item3, p.item4, p.item5, p.item6, p.summoner1_id,
			p.summoner2_id, p.lane, p.participant_id, p.puuid,
			p.riot_id_game_name, p.riot_id_tagline,
			p.total_damage_dealt_to_champions, p.total_damage_taken, p.win
		FROM lol_matches m
		LEFT JOIN participants p ON m.id = p.match_id
		WHERE m.id IN (%s)
		ORDER BY m.started_at DESC, p.participant_id ASC
	`, strings.Join(placeholders, ", "))

	// Step 4: Execute second query
	fullRows, err := db.SQL.Query(fullQuery, fullArgs...)
	if err != nil {
		logging.Error("Failed to query full match data from database", "error", err)
		return nil, err
	}
	defer fullRows.Close()

	// Step 5: Scan and group results by match ID
	matchMap := make(map[int64]*models.LeagueOfLegendsMatch)
	var orderedMatchIDs []int64 // Preserve order

	for fullRows.Next() {
		var summary models.LeagueOfLegendsMatchSummary

		// Use nullable types for participant fields to handle LEFT JOIN NULLs
		var (
			nullMatchID                     sql.NullInt64
			nullChampionID                  sql.NullInt64
			nullKills                       sql.NullInt64
			nullDeaths                      sql.NullInt64
			nullAssists                     sql.NullInt64
			nullTotalMinionsKilled          sql.NullInt64
			nullDoubleKills                 sql.NullInt64
			nullTripleKills                 sql.NullInt64
			nullQuadraKills                 sql.NullInt64
			nullPentaKills                  sql.NullInt64
			nullItem0                       sql.NullInt64
			nullItem1                       sql.NullInt64
			nullItem2                       sql.NullInt64
			nullItem3                       sql.NullInt64
			nullItem4                       sql.NullInt64
			nullItem5                       sql.NullInt64
			nullItem6                       sql.NullInt64
			nullSummoner1ID                 sql.NullInt64
			nullSummoner2ID                 sql.NullInt64
			nullLane                        sql.NullString
			nullParticipantID               sql.NullInt64
			nullPUUID                       sql.NullString
			nullRiotIDGameName              sql.NullString
			nullRiotIDTagline               sql.NullString
			nullTotalDamageDealtToChampions sql.NullInt64
			nullTotalDamageTaken            sql.NullInt64
			nullWin                         sql.NullBool
		)

		err := fullRows.Scan(
			&summary.ID, &summary.StartedAt, &summary.Duration, &summary.ReplaySynced,
			&nullMatchID, &nullChampionID, &nullKills,
			&nullDeaths, &nullAssists, &nullTotalMinionsKilled,
			&nullDoubleKills, &nullTripleKills, &nullQuadraKills,
			&nullPentaKills, &nullItem0, &nullItem1, &nullItem2,
			&nullItem3, &nullItem4, &nullItem5, &nullItem6,
			&nullSummoner1ID, &nullSummoner2ID, &nullLane,
			&nullParticipantID, &nullPUUID, &nullRiotIDGameName,
			&nullRiotIDTagline, &nullTotalDamageDealtToChampions,
			&nullTotalDamageTaken, &nullWin,
		)
		if err != nil {
			logging.Error("Failed to scan full match data row", "error", err)
			return nil, err
		}

		// Check if match already exists in map
		if _, exists := matchMap[summary.ID]; !exists {
			matchMap[summary.ID] = &models.LeagueOfLegendsMatch{
				Summary:      summary,
				Participants: []models.LeagueOfLegendsMatchParticipantSummary{},
			}
			orderedMatchIDs = append(orderedMatchIDs, summary.ID)
		}

		// Only append participant if data is present (not NULL from LEFT JOIN)
		// This handles matches that exist without participants (orphaned matches)
		if nullMatchID.Valid {
			participant := models.LeagueOfLegendsMatchParticipantSummary{
				GameID:                      nullMatchID.Int64,
				ChampionID:                  int(nullChampionID.Int64),
				Kills:                       int(nullKills.Int64),
				Deaths:                      int(nullDeaths.Int64),
				Assists:                     int(nullAssists.Int64),
				TotalMinionsKilled:          int(nullTotalMinionsKilled.Int64),
				DoubleKills:                 int(nullDoubleKills.Int64),
				TripleKills:                 int(nullTripleKills.Int64),
				QuadraKills:                 int(nullQuadraKills.Int64),
				PentaKills:                  int(nullPentaKills.Int64),
				Item0:                       int(nullItem0.Int64),
				Item1:                       int(nullItem1.Int64),
				Item2:                       int(nullItem2.Int64),
				Item3:                       int(nullItem3.Int64),
				Item4:                       int(nullItem4.Int64),
				Item5:                       int(nullItem5.Int64),
				Item6:                       int(nullItem6.Int64),
				Summoner1ID:                 int(nullSummoner1ID.Int64),
				Summoner2ID:                 int(nullSummoner2ID.Int64),
				Lane:                        nullLane.String,
				ParticipantID:               int(nullParticipantID.Int64),
				PUUID:                       nullPUUID.String,
				RiotIDGameName:              nullRiotIDGameName.String,
				RiotIDTagline:               nullRiotIDTagline.String,
				TotalDamageDealtToChampions: int(nullTotalDamageDealtToChampions.Int64),
				TotalDamageTaken:            int(nullTotalDamageTaken.Int64),
				Win:                         nullWin.Bool,
			}

			// Append participant
			matchMap[summary.ID].Participants = append(
				matchMap[summary.ID].Participants,
				participant,
			)
		}
	}
	if err := fullRows.Err(); err != nil {
		logging.Error("Error iterating over full match data rows", "error", err)
		return nil, err
	}

	// Step 6: Convert map to ordered slice
	result := make([]models.LeagueOfLegendsMatch, 0, len(orderedMatchIDs))
	for _, id := range orderedMatchIDs {
		result = append(result, *matchMap[id])
	}

	return result, nil
}
