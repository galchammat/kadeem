CREATE TABLE IF NOT EXISTS league_of_legends_accounts (
    puuid     TEXT    NOT NULL PRIMARY KEY,
    tag_line  TEXT,
    game_name TEXT,
    region    TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);