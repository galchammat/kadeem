ALTER TABLE lol_matches DROP COLUMN queue_id;
ALTER TABLE participants DROP COLUMN champ_level;
DROP TABLE IF EXISTS player_ranks;
DROP INDEX IF EXISTS idx_player_ranks_puuid_time;
DROP INDEX IF EXISTS idx_player_ranks_lookup;
