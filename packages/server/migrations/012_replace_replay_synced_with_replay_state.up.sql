ALTER TABLE lol_matches ADD COLUMN replay_s3_key TEXT;
ALTER TABLE lol_matches ADD COLUMN replay_sync_error TEXT;
ALTER TABLE lol_matches ADD COLUMN replay_sync_attempted_at TIMESTAMPTZ;

UPDATE lol_matches
SET replay_s3_key = CONCAT('legacy/lol/replays/', id, '.rofl')
WHERE replay_synced = TRUE;

ALTER TABLE lol_matches DROP COLUMN replay_synced;

CREATE INDEX IF NOT EXISTS idx_lol_matches_replay_pending
ON lol_matches(id DESC)
WHERE replay_s3_key IS NULL;
