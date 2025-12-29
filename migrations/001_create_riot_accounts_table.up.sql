CREATE TABLE
	IF NOT EXISTS league_of_legends_accounts (
		puuid VARCHAR(78) NOT NULL PRIMARY KEY,
		streamer_id INTEGER NOT NULL REFERENCES streamers (id) ON DELETE CASCADE,
		tag_line VARCHAR(5),
		game_name VARCHAR(16),
		region VARCHAR(4)
	);