CREATE TABLE IF NOT EXISTS stream_events (
    id          BIGSERIAL PRIMARY KEY,
    channel_id  VARCHAR(30) NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    event_type  VARCHAR(30) NOT NULL,
    title       TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    timestamp   BIGINT NOT NULL,
    value       TEXT,
    external_id TEXT
);

-- Prevents duplicate events from repeated syncs
CREATE UNIQUE INDEX IF NOT EXISTS idx_stream_events_external
    ON stream_events(channel_id, external_id)
    WHERE external_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_stream_events_channel_time
    ON stream_events(channel_id, timestamp DESC);
