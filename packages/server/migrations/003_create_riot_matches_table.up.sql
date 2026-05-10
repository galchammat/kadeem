CREATE TABLE IF NOT EXISTS lol_matches (
	id BIGINT PRIMARY KEY,
  region VARCHAR(5),
	started_at BIGINT,
	duration INTEGER,
	queue_id INTEGER,
	status TEXT NOT NULL DEFAULT 'pending',
  updated_at TIMESTAMPTZ,
	replay_uri TEXT,
	replay_status TEXT NOT NULL DEFAULT 'pending',
  replay_updated_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_lol_matches_started_at ON lol_matches(started_at DESC);
CREATE UNIQUE INDEX IF NOT EXISTS idx_lol_matches_id ON lol_matches(id, region)
