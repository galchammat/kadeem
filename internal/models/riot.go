package models

type LeagueOfLegendsAccount struct {
	PUUID    string `json:"puuid"`
	TagLine  string `json:"tagLine"`
	GameName string `json:"gameName"`
	Region   string
}
