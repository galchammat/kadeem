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
