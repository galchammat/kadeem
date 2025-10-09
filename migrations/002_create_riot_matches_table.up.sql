CREATE TABLE IF NOT EXISTS league_of_legends_matches (
    match_id varchar(16) PRIMARY KEY,
    account_id varchar(78) NOT NULL REFERENCES league_of_legends_accounts(puuid) ON DELETE CASCADE,
    champion_id integer NOT NULL,
    kills integer NOT NULL,
    deaths integer NOT NULL,
    assists integer NOT NULL,
    win boolean NOT NULL,
    game_creation bigint NOT NULL,
    game_duration integer NOT NULL,
    damage integer NOT NULL
)