package database

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/galchammat/kadeem/internal/model"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) (*DB, func()) {
	// Create a temporary database file
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	sqlDB, err := sql.Open("sqlite3", "file://"+dbPath)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Enable foreign keys
	if _, err := sqlDB.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		sqlDB.Close()
		t.Fatalf("Failed to enable foreign keys: %v", err)
	}

	db := &DB{SQL: sqlDB}

	// Create tables
	createMatchesTable := `
		CREATE TABLE IF NOT EXISTS lol_matches (
			id BIGINT PRIMARY KEY,
			started_at BIGINT NOT NULL,
			duration INTEGER NOT NULL,
			queue_id INTEGER NOT NULL DEFAULT 0,
			replay_synced BOOLEAN DEFAULT FALSE
		);
	`
	if _, err := sqlDB.Exec(createMatchesTable); err != nil {
		sqlDB.Close()
		t.Fatalf("Failed to create lol_matches table: %v", err)
	}

	createParticipantsTable := `
		CREATE TABLE IF NOT EXISTS participants (
			match_id BIGINT NOT NULL REFERENCES lol_matches(id) ON DELETE CASCADE,
			champion_id INTEGER NOT NULL,
			champ_level INTEGER NOT NULL DEFAULT 1,
			kills INTEGER NOT NULL,
			deaths INTEGER NOT NULL,
			assists INTEGER NOT NULL,
			total_minions_killed INTEGER NOT NULL,
			double_kills INTEGER NOT NULL,
			triple_kills INTEGER NOT NULL,
			quadra_kills INTEGER NOT NULL,
			penta_kills INTEGER NOT NULL,
			item0 INTEGER NOT NULL,
			item1 INTEGER NOT NULL,
			item2 INTEGER NOT NULL,
			item3 INTEGER NOT NULL,
			item4 INTEGER NOT NULL,
			item5 INTEGER NOT NULL,
			item6 INTEGER NOT NULL,
			summoner1_id INTEGER NOT NULL,
			summoner2_id INTEGER NOT NULL,
			lane TEXT NOT NULL,
			participant_id INTEGER NOT NULL,
			puuid VARCHAR(78) NOT NULL,
			riot_id_game_name TEXT NOT NULL,
			riot_id_tagline TEXT NOT NULL,
			total_damage_dealt_to_champions INTEGER NOT NULL,
			total_damage_taken INTEGER NOT NULL,
			win BOOLEAN NOT NULL,
			PRIMARY KEY (match_id, participant_id)
		);
	`
	if _, err := sqlDB.Exec(createParticipantsTable); err != nil {
		sqlDB.Close()
		t.Fatalf("Failed to create participants table: %v", err)
	}

	cleanup := func() {
		sqlDB.Close()
		os.RemoveAll(tmpDir)
	}

	return db, cleanup
}

func TestListLolMatches_HandlesOrphanedMatches(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Insert a match without participants (orphaned match)
	matchID := int64(12345)
	startedAt := int64(1600000000)
	duration := 1800

	_, err := db.SQL.Exec(
		"INSERT INTO lol_matches (id, started_at, duration, replay_synced) VALUES (?, ?, ?, ?)",
		matchID, startedAt, duration, false,
	)
	if err != nil {
		t.Fatalf("Failed to insert orphaned match: %v", err)
	}

	// Try to list matches - this should NOT crash
	limit := 10
	offset := 0
	matches, err := db.ListLolMatches(&model.LolMatchFilter{}, limit, offset)

	if err != nil {
		t.Fatalf("ListLolMatches failed with orphaned match: %v", err)
	}

	// Should return 1 match
	if len(matches) != 1 {
		t.Errorf("Expected 1 match, got %d", len(matches))
	}

	// Match should have empty participants array
	if len(matches[0].Participants) != 0 {
		t.Errorf("Expected 0 participants for orphaned match, got %d", len(matches[0].Participants))
	}

	// Verify match summary data
	if matches[0].Summary.ID != matchID {
		t.Errorf("Expected match ID %d, got %d", matchID, matches[0].Summary.ID)
	}
	if *matches[0].Summary.StartedAt != startedAt {
		t.Errorf("Expected startedAt %d, got %d", startedAt, *matches[0].Summary.StartedAt)
	}
	if *matches[0].Summary.Duration != duration {
		t.Errorf("Expected duration %d, got %d", duration, *matches[0].Summary.Duration)
	}

	t.Logf("Successfully handled orphaned match: ID=%d, Participants=%d",
		matches[0].Summary.ID, len(matches[0].Participants))
}

func TestInsertLolMatchWithParticipants_Transaction(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	matchID := int64(67890)
	startedAt := int64(1600000000)
	duration := 2000

	summary := &model.LolMatchSummary{
		ID:        matchID,
		StartedAt: &startedAt,
		Duration:  &duration,
	}

	participants := []model.LolMatchParticipantSummary{
		{
			GameID:                      matchID,
			ChampionID:                  1,
			Kills:                       10,
			Deaths:                      2,
			Assists:                     5,
			TotalMinionsKilled:          150,
			DoubleKills:                 1,
			TripleKills:                 0,
			QuadraKills:                 0,
			PentaKills:                  0,
			Item0:                       1001,
			Item1:                       1002,
			Item2:                       1003,
			Item3:                       1004,
			Item4:                       1005,
			Item5:                       1006,
			Item6:                       0,
			Summoner1ID:                 4,
			Summoner2ID:                 7,
			Lane:                        "MIDDLE",
			ParticipantID:               1,
			PUUID:                       "test-puuid-1",
			RiotIDGameName:              "TestPlayer",
			RiotIDTagline:               "NA1",
			TotalDamageDealtToChampions: 20000,
			TotalDamageTaken:            15000,
			Win:                         true,
		},
	}

	// Test successful transaction
	err := db.InsertLolMatchWithParticipants(summary, participants)
	if err != nil {
		t.Fatalf("InsertLolMatchWithParticipants failed: %v", err)
	}

	// Verify match was inserted
	var count int
	err = db.SQL.QueryRow("SELECT COUNT(*) FROM lol_matches WHERE id = ?", matchID).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query matches: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 match, got %d", count)
	}

	// Verify participant was inserted
	err = db.SQL.QueryRow("SELECT COUNT(*) FROM participants WHERE match_id = ?", matchID).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query participants: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 participant, got %d", count)
	}

	t.Logf("Successfully inserted match with participants using transaction")
}

func TestListLolMatches_WithCompleteMatch(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	matchID := int64(11111)
	startedAt := int64(1600000000)
	duration := 2500

	summary := &model.LolMatchSummary{
		ID:        matchID,
		StartedAt: &startedAt,
		Duration:  &duration,
	}

	participants := []model.LolMatchParticipantSummary{
		{
			GameID:                      matchID,
			ChampionID:                  10,
			Kills:                       5,
			Deaths:                      3,
			Assists:                     8,
			TotalMinionsKilled:          120,
			DoubleKills:                 0,
			TripleKills:                 0,
			QuadraKills:                 0,
			PentaKills:                  0,
			Item0:                       2001,
			Item1:                       2002,
			Item2:                       2003,
			Item3:                       2004,
			Item4:                       2005,
			Item5:                       2006,
			Item6:                       0,
			Summoner1ID:                 4,
			Summoner2ID:                 7,
			Lane:                        "TOP",
			ParticipantID:               1,
			PUUID:                       "test-puuid-complete",
			RiotIDGameName:              "CompletePlayer",
			RiotIDTagline:               "EUW",
			TotalDamageDealtToChampions: 18000,
			TotalDamageTaken:            12000,
			Win:                         false,
		},
		{
			GameID:                      matchID,
			ChampionID:                  20,
			Kills:                       8,
			Deaths:                      4,
			Assists:                     12,
			TotalMinionsKilled:          180,
			DoubleKills:                 2,
			TripleKills:                 0,
			QuadraKills:                 0,
			PentaKills:                  0,
			Item0:                       3001,
			Item1:                       3002,
			Item2:                       3003,
			Item3:                       3004,
			Item4:                       3005,
			Item5:                       3006,
			Item6:                       0,
			Summoner1ID:                 11,
			Summoner2ID:                 14,
			Lane:                        "JUNGLE",
			ParticipantID:               2,
			PUUID:                       "test-puuid-complete-2",
			RiotIDGameName:              "CompletePlayer2",
			RiotIDTagline:               "EUW",
			TotalDamageDealtToChampions: 22000,
			TotalDamageTaken:            14000,
			Win:                         false,
		},
	}

	// Insert complete match
	err := db.InsertLolMatchWithParticipants(summary, participants)
	if err != nil {
		t.Fatalf("Failed to insert complete match: %v", err)
	}

	// List matches
	limit := 10
	offset := 0
	matches, err := db.ListLolMatches(&model.LolMatchFilter{}, limit, offset)

	if err != nil {
		t.Fatalf("ListLolMatches failed: %v", err)
	}

	if len(matches) != 1 {
		t.Fatalf("Expected 1 match, got %d", len(matches))
	}

	// Should have 2 participants
	if len(matches[0].Participants) != 2 {
		t.Errorf("Expected 2 participants, got %d", len(matches[0].Participants))
	}

	// Verify participant data
	for i, p := range matches[0].Participants {
		if p.GameID != matchID {
			t.Errorf("Participant %d: expected GameID %d, got %d", i, matchID, p.GameID)
		}
		if p.ParticipantID != i+1 {
			t.Errorf("Participant %d: expected ParticipantID %d, got %d", i, i+1, p.ParticipantID)
		}
	}

	t.Logf("Successfully listed complete match with %d participants", len(matches[0].Participants))
}
