CREATE TABLE IF NOT EXISTS lol_matches (
	id BIGINT PRIMARY KEY,
	started_at BIGINT,
	duration INTEGER,
	queue_id INTEGER,
	status TEXT NOT NULL DEFAULT 'pending',
	replay_status TEXT NOT NULL DEFAULT 'pending',
	replay_uri TEXT,
	replay_sync_error TEXT,
	replay_sync_attempted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_lol_matches_started_at ON lol_matches(started_at DESC);
CREATE INDEX IF NOT EXISTS idx_lol_matches_replay_pending
ON lol_matches(id DESC)
WHERE replay_status = 'pending';
