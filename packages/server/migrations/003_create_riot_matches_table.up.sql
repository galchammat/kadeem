CREATE TABLE IF NOT EXISTS lol_matches (
	id BIGINT PRIMARY KEY,
	started_at BIGINT NOT NULL,
	duration INTEGER NOT NULL,
	replay_synced BOOLEAN DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_lol_matches_started_at ON lol_matches(started_at DESC);
