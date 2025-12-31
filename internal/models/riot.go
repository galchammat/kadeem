package models

type LeagueOfLegendsAccount struct {
	PUUID      string `json:"puuid" db:"puuid,primarykey"`
	StreamerID int    `json:"streamerId,omitempty" db:"streamer_id"`
	TagLine    string `json:"tagLine" db:"tag_line"`
	GameName   string `json:"gameName" db:"game_name"`
	Region     string `json:"region,omitempty" db:"region"`
	SyncedAt   *int64 `json:"syncedAt" db:"synced_at"`
}

type LeagueOfLegendsMatchSummary struct {
	AccountID string `json:"accountId" db:"account_id,primarykey"`
}
