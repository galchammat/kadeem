package models

type LeagueOfLegendsAccount struct {
	PUUID    string `json:"puuid" db:"puuid,primarykey"`
	TagLine  string `json:"tagLine" db:"tag_line"`
	GameName string `json:"gameName" db:"game_name"`
	Region   string `json:"region,omitempty" db:"region"`
	Streamer string `json:"streamer,omitempty" db:"streamer"`
}
