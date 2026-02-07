package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/galchammat/kadeem/internal/model"
)

// SQLExecutor interface allows functions to accept either *sql.DB or *sql.Tx
type SQLExecutor interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

func (db *DB) InsertLolMatchSummary(summary *model.LeagueOfLegendsMatchSummary) error {
	return db.insertLolMatchSummaryExec(db.SQL, summary)
}

func (db *DB) insertLolMatchSummaryExec(exec SQLExecutor, summary *model.LeagueOfLegendsMatchSummary) error {
	if summary == nil || summary.StartedAt == nil || summary.Duration == nil {
		return fmt.Errorf("match summary is missing required fields")
	}

	replaySynced := false
	if summary.ReplaySynced != nil {
		replaySynced = *summary.ReplaySynced
	}

	queueId := 0
	if summary.QueueId != nil {
		queueId = *summary.QueueId
	}

	query := `
		INSERT INTO lol_matches
		(id, started_at, duration, queue_id, replay_synced)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			started_at = EXCLUDED.started_at,
			duration = EXCLUDED.duration,
			queue_id = EXCLUDED.queue_id,
			replay_synced = EXCLUDED.replay_synced`

	_, err := exec.Exec(query, summary.ID, *summary.StartedAt, *summary.Duration, queueId, replaySynced)
	if err != nil {
		logging.Error("Failed to insert match summary into database", "matchID", summary.ID, "error", err)
	}
	return err
}

// allowedMatchColumns is the set of columns that can be updated via UpdateLolMatch.
var allowedMatchColumns = map[string]bool{
	"started_at":    true,
	"duration":      true,
	"queue_id":      true,
	"replay_synced": true,
}

func (db *DB) UpdateLolMatch(matchID int64, updates map[string]any) (bool, error) {
	var setClauses []string
	var args []any
	argN := 1

	for column, value := range updates {
		if !allowedMatchColumns[column] {
			return false, fmt.Errorf("disallowed column: %s", column)
		}
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, argN))
		args = append(args, value)
		argN++
	}
	if len(setClauses) == 0 {
		return false, nil
	}
	args = append(args, matchID)

	query := `UPDATE lol_matches SET ` + strings.Join(setClauses, ", ") + fmt.Sprintf(` WHERE id = $%d`, argN)

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
	summary *model.LeagueOfLegendsMatchSummary,
	participants []model.LeagueOfLegendsMatchParticipantSummary,
) error {
	tx, err := db.SQL.Begin()
	if err != nil {
		logging.Error("Failed to begin transaction for match insert", "matchID", summary.ID, "error", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck // rollback after commit is a no-op

	if err := db.insertLolMatchSummaryExec(tx, summary); err != nil {
		return err
	}

	for i, participant := range participants {
		participant.GameID = summary.ID
		if err := db.insertLolMatchParticipantSummaryExec(tx, &participant); err != nil {
			logging.Error(
				"Failed to insert participant, rolling back",
				"matchID", summary.ID,
				"participantIndex", i,
				"participantID", participant.ParticipantID,
				"error", err,
			)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		logging.Error("Failed to commit match transaction", "matchID", summary.ID, "error", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logging.Debug("Successfully inserted match with participants", "matchID", summary.ID, "participantCount", len(participants))
	return nil
}

func (db *DB) InsertLolMatchParticipantSummary(participant *model.LeagueOfLegendsMatchParticipantSummary) error {
	return db.insertLolMatchParticipantSummaryExec(db.SQL, participant)
}

func (db *DB) insertLolMatchParticipantSummaryExec(exec SQLExecutor, participant *model.LeagueOfLegendsMatchParticipantSummary) error {
	if participant == nil {
		return fmt.Errorf("participant summary cannot be nil")
	}

	query := `
		INSERT INTO participants (
			match_id, champion_id, champ_level, kills, deaths, assists,
			total_minions_killed, double_kills, triple_kills, quadra_kills, penta_kills,
			item0, item1, item2, item3, item4, item5, item6,
			summoner1_id, summoner2_id, lane, participant_id, puuid,
			riot_id_game_name, riot_id_tagline,
			total_damage_dealt_to_champions, total_damage_taken, win
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28)
		ON CONFLICT (match_id, participant_id) DO UPDATE SET
			champion_id = EXCLUDED.champion_id,
			champ_level = EXCLUDED.champ_level,
			kills = EXCLUDED.kills,
			deaths = EXCLUDED.deaths,
			assists = EXCLUDED.assists,
			total_minions_killed = EXCLUDED.total_minions_killed,
			double_kills = EXCLUDED.double_kills,
			triple_kills = EXCLUDED.triple_kills,
			quadra_kills = EXCLUDED.quadra_kills,
			penta_kills = EXCLUDED.penta_kills,
			item0 = EXCLUDED.item0,
			item1 = EXCLUDED.item1,
			item2 = EXCLUDED.item2,
			item3 = EXCLUDED.item3,
			item4 = EXCLUDED.item4,
			item5 = EXCLUDED.item5,
			item6 = EXCLUDED.item6,
			summoner1_id = EXCLUDED.summoner1_id,
			summoner2_id = EXCLUDED.summoner2_id,
			lane = EXCLUDED.lane,
			puuid = EXCLUDED.puuid,
			riot_id_game_name = EXCLUDED.riot_id_game_name,
			riot_id_tagline = EXCLUDED.riot_id_tagline,
			total_damage_dealt_to_champions = EXCLUDED.total_damage_dealt_to_champions,
			total_damage_taken = EXCLUDED.total_damage_taken,
			win = EXCLUDED.win`

	_, err := exec.Exec(
		query,
		participant.GameID,
		participant.ChampionID,
		participant.ChampLevel,
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

func (db *DB) ListLolMatches(filter *model.LolMatchFilter, limit int, offset int) ([]model.LeagueOfLegendsMatch, error) {
	// Clamp limit
	if limit <= 0 || limit > 100 {
		limit = 100
	}

	// Step 1: Build query for matching match IDs
	matchIDQuery := `
		SELECT DISTINCT m.id
		FROM lol_matches m
		LEFT JOIN participants p ON m.id = p.match_id`

	var where []string
	var args []any
	argN := 1

	if filter != nil {
		if filter.PUUID != nil {
			where = append(where, fmt.Sprintf("p.puuid = $%d", argN))
			args = append(args, *filter.PUUID)
			argN++
		}
		if filter.MatchID != nil {
			where = append(where, fmt.Sprintf("m.id = $%d", argN))
			args = append(args, *filter.MatchID)
			argN++
		}
		if filter.StartedAtMin != nil {
			where = append(where, fmt.Sprintf("m.started_at >= $%d", argN))
			args = append(args, *filter.StartedAtMin)
			argN++
		}
		if filter.StartedAtMax != nil {
			where = append(where, fmt.Sprintf("m.started_at <= $%d", argN))
			args = append(args, *filter.StartedAtMax)
			argN++
		}
		if filter.ReplaySynced != nil {
			where = append(where, fmt.Sprintf("m.replay_synced = $%d", argN))
			args = append(args, *filter.ReplaySynced)
			argN++
		}
		if filter.ChampionID != nil {
			where = append(where, fmt.Sprintf("p.champion_id = $%d", argN))
			args = append(args, *filter.ChampionID)
			argN++
		}
		if filter.Lane != nil {
			where = append(where, fmt.Sprintf("p.lane = $%d", argN))
			args = append(args, *filter.Lane)
			argN++
		}
		if filter.Win != nil {
			where = append(where, fmt.Sprintf("p.win = $%d", argN))
			args = append(args, *filter.Win)
			argN++
		}
	}

	if len(where) > 0 {
		matchIDQuery += " WHERE " + strings.Join(where, " AND ")
	}

	matchIDQuery += fmt.Sprintf(" ORDER BY m.id DESC LIMIT $%d OFFSET $%d", argN, argN+1)
	args = append(args, limit, offset)

	// Step 2: Execute to get match IDs
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
			return nil, err
		}
		matchIDs = append(matchIDs, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(matchIDs) == 0 {
		return []model.LeagueOfLegendsMatch{}, nil
	}

	// Step 3: Fetch full match data for the selected IDs
	placeholders := make([]string, len(matchIDs))
	fullArgs := make([]any, len(matchIDs))
	for i, id := range matchIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		fullArgs[i] = id
	}

	fullQuery := fmt.Sprintf(`
		SELECT 
			m.id, m.started_at, m.duration, m.queue_id, m.replay_synced,
			p.match_id, p.champion_id, p.champ_level, p.kills, p.deaths, p.assists,
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

	// Step 4: Execute
	fullRows, err := db.SQL.Query(fullQuery, fullArgs...)
	if err != nil {
		logging.Error("Failed to query full match data from database", "error", err)
		return nil, err
	}
	defer fullRows.Close()

	// Step 5: Scan and group by match ID
	matchMap := make(map[int64]*model.LeagueOfLegendsMatch)
	var orderedMatchIDs []int64

	for fullRows.Next() {
		var summary model.LeagueOfLegendsMatchSummary

		var (
			nullMatchID                     sql.NullInt64
			nullChampionID                  sql.NullInt64
			nullChampLevel                  sql.NullInt64
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
			nullQueueId                     sql.NullInt64
		)

		err := fullRows.Scan(
			&summary.ID, &summary.StartedAt, &summary.Duration, &nullQueueId, &summary.ReplaySynced,
			&nullMatchID, &nullChampionID, &nullChampLevel, &nullKills,
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

		if _, exists := matchMap[summary.ID]; !exists {
			if nullQueueId.Valid {
				qid := int(nullQueueId.Int64)
				summary.QueueId = &qid
			}
			matchMap[summary.ID] = &model.LeagueOfLegendsMatch{
				Summary:      summary,
				Participants: []model.LeagueOfLegendsMatchParticipantSummary{},
			}
			orderedMatchIDs = append(orderedMatchIDs, summary.ID)
		}

		if nullMatchID.Valid {
			participant := model.LeagueOfLegendsMatchParticipantSummary{
				GameID:                      nullMatchID.Int64,
				ChampionID:                  int(nullChampionID.Int64),
				ChampLevel:                  int(nullChampLevel.Int64),
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

	// Step 6: Convert to ordered slice
	result := make([]model.LeagueOfLegendsMatch, 0, len(orderedMatchIDs))
	for _, id := range orderedMatchIDs {
		result = append(result, *matchMap[id])
	}

	return result, nil
}
