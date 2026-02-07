-- Add queue_id to lol_matches
ALTER TABLE lol_matches ADD COLUMN queue_id INTEGER;

-- Add champ_level to participants
ALTER TABLE participants ADD COLUMN champ_level INTEGER NOT NULL DEFAULT 1;

-- Create rank history table
CREATE TABLE IF NOT EXISTS player_ranks (
    puuid VARCHAR(78) NOT NULL,
    timestamp BIGINT NOT NULL,
    tier TEXT NOT NULL,
    rank TEXT NOT NULL,
    league_points INTEGER NOT NULL,
    wins INTEGER NOT NULL,
    losses INTEGER NOT NULL,
    queue_id INTEGER NOT NULL,
    PRIMARY KEY (puuid, timestamp, queue_id)
);

CREATE INDEX idx_player_ranks_puuid_time ON player_ranks(puuid, timestamp DESC);
CREATE INDEX idx_player_ranks_lookup ON player_ranks(puuid, queue_id, timestamp DESC);
