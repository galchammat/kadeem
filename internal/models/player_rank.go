package models

type PlayerRank struct {
	PUUID        string `json:"puuid" db:"puuid"`
	Timestamp    int64  `json:"timestamp" db:"timestamp"`
	Tier         string `json:"tier" db:"tier"`
	Rank         string `json:"rank" db:"rank"`
	LeaguePoints int    `json:"leaguePoints" db:"league_points"`
	Wins         int    `json:"wins" db:"wins"`
	Losses       int    `json:"losses" db:"losses"`
	QueueId      int    `json:"queueId" db:"queue_id"`
}
