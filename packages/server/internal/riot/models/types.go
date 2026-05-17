package models

import "time"

type APIReplaysResponse struct {
	URLs []string `json:"matchFileURLs"`
}

type Match struct {
	Summary      MatchSummary              `json:"summary" db:"-"`
	Participants []MatchParticipantSummary `json:"participants" db:"-"`
}

type MatchSummary struct {
	ID              int64      `json:"gameId" db:"match_id"`
	Region          string     `json:"region" db:"region"`
	StartedAt       int64      `json:"startedAt" db:"started_at"`
	Duration        int        `json:"duration" db:"duration"`
	QueueID         int        `json:"queueId" db:"queue_id"`
	Status          string     `json:"status" db:"status"`
	UpdatedAt       *time.Time `json:"updatedAt" db:"updated_at"`
	ReplayStatus    string     `json:"replayStatus" db:"replay_status"`
	ReplayURI       *string    `json:"replayUri,omitempty" db:"replay_uri"`
	ReplayUpdatedAt *time.Time `json:"replayUpdatedAt,omitempty" db:"replay_updated_at"`
}

type MatchParticipantSummary struct {
	GameID                      int64  `json:"gameId" db:"match_id"`
	ChampionID                  int    `json:"championId" db:"champion_id"`
	ChampLevel                  int    `json:"champLevel" db:"champ_level"`
	Kills                       int    `json:"kills" db:"kills"`
	Deaths                      int    `json:"deaths" db:"deaths"`
	Assists                     int    `json:"assists" db:"assists"`
	TotalMinionsKilled          int    `json:"totalMinionsKilled" db:"total_minions_killed"`
	DoubleKills                 int    `json:"doubleKills" db:"double_kills"`
	TripleKills                 int    `json:"tripleKills" db:"triple_kills"`
	QuadraKills                 int    `json:"quadraKills" db:"quadra_kills"`
	PentaKills                  int    `json:"pentaKills" db:"penta_kills"`
	Item0                       int    `json:"item0" db:"item0"`
	Item1                       int    `json:"item1" db:"item1"`
	Item2                       int    `json:"item2" db:"item2"`
	Item3                       int    `json:"item3" db:"item3"`
	Item4                       int    `json:"item4" db:"item4"`
	Item5                       int    `json:"item5" db:"item5"`
	Item6                       int    `json:"item6" db:"item6"`
	Summoner1ID                 int    `json:"summoner1Id" db:"summoner1_id"`
	Summoner2ID                 int    `json:"summoner2Id" db:"summoner2_id"`
	Lane                        string `json:"lane" db:"lane"`
	ParticipantID               int    `json:"participantId" db:"participant_id"`
	PUUID                       string `json:"puuid" db:"puuid"`
	RiotIDGameName              string `json:"riotIdGameName" db:"riot_id_game_name"`
	RiotIDTagline               string `json:"riotIdTagline" db:"riot_id_tagline"`
	TotalDamageDealtToChampions int    `json:"totalDamageDealtToChampions" db:"total_damage_dealt_to_champions"`
	TotalDamageTaken            int    `json:"totalDamageTaken" db:"total_damage_taken"`
	Win                         bool   `json:"win" db:"win"`
}

type MatchFilter struct {
	MatchID      *int64
	StartedAtMin *int64
	StartedAtMax *int64
	HasReplay    *bool
	PUUID        *string
	ChampionID   *int
	Lane         *string
	Win          *bool
}

type Account struct {
	PUUID      string `json:"puuid" db:"puuid,primarykey"`
	StreamerID int    `json:"streamerId,omitempty" db:"streamer_id"`
	TagLine    string `json:"tagLine" db:"tag_line"`
	GameName   string `json:"gameName" db:"game_name"`
	Region     string `json:"region,omitempty" db:"region"`
	SyncedAt   *int64 `json:"syncedAt" db:"synced_at"`
}

type PlayerRank struct {
	PUUID        string `json:"puuid" db:"puuid"`
	Timestamp    int64  `json:"timestamp" db:"timestamp"`
	Tier         string `json:"tier" db:"tier"`
	Rank         string `json:"rank" db:"rank"`
	LeaguePoints int    `json:"leaguePoints" db:"league_points"`
	Wins         int    `json:"wins" db:"wins"`
	Losses       int    `json:"losses" db:"losses"`
	QueueID      int    `json:"queueId" db:"queue_id"`
}

type MatchDetails struct {
	Info struct {
		ID           int64                     `json:"gameId"`
		Region       string                    `json:"region"`
		QueueID      int                       `json:"queueId"`
		StartedAt    int64                     `json:"gameStartTimestamp"`
		Duration     int                       `json:"gameDuration"`
		Participants []MatchParticipantSummary `json:"participants"`
	} `json:"info"`
}
