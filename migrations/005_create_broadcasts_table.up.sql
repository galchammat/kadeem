CREATE TABLE IF NOT EXISTS broadcasts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    channel_id VARCHAR(30) REFERENCES streams(id) ON DELETE CASCADE,
    title VARCHAR(255),
    url VARCHAR(255) NOT NULL,
    thumbnail_url VARCHAR(255) NOT NULL,
    viewable VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    published_at TIMESTAMP NOT NULL,
    duration INTEGER,
);
