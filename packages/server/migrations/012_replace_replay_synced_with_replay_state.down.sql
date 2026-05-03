DROP INDEX IF EXISTS idx_lol_matches_replay_pending;

ALTER TABLE lol_matches ADD COLUMN replay_synced BOOLEAN DEFAULT FALSE;

UPDATE lol_matches
SET replay_synced = replay_s3_key IS NOT NULL;

ALTER TABLE lol_matches DROP COLUMN replay_sync_attempted_at;
ALTER TABLE lol_matches DROP COLUMN replay_sync_error;
ALTER TABLE lol_matches DROP COLUMN replay_s3_key;
