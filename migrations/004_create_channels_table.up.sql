CREATE TABLE IF NOT EXISTS channels (
    id VARCHAR(30) PRIMARY KEY,
    streamer_id INTEGER REFERENCES streamers(id) ON DELETE CASCADE,
    platform VARCHAR(10) NOT NULL,
    channel_name VARCHAR(30) NOT NULL,
    avatar_url VARCHAR(255) NOT NULL,
    synced_at TIMESTAMP DEFAULT NULL
);