CREATE TABLE IF NOT EXISTS streams (
    id SERIAL PRIMARY KEY,
    streamer_id INTEGER REFERENCES streamers(id) ON DELETE CASCADE,
    platform VARCHAR(10) NOT NULL,
    channel_name VARCHAR(30) NOT NULL,
    channel_id VARCHAR(30) NOT NULL
)