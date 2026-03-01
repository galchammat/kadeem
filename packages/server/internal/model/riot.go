package model

type LolApiReplaysResponse struct {
	URLs []string `json:"matchFileURLs"`
}

type LolAccount struct {
	PUUID      string `json:"puuid" db:"puuid,primarykey"`
	StreamerID int    `json:"streamerId,omitempty" db:"streamer_id"`
	TagLine    string `json:"tagLine" db:"tag_line"`
	GameName   string `json:"gameName" db:"game_name"`
	Region     string `json:"region,omitempty" db:"region"`
	SyncedAt   *int64 `json:"syncedAt" db:"synced_at"`
}

type LolMatch struct {
	Summary      LolMatchSummary              `json:"summary" db:"-"`
	Participants []LolMatchParticipantSummary `json:"participants" db:"-"`
	ReplayURL    *string                      `json:"replay,omitempty" db:"replay"`
}

type LolMatchSummary struct {
	ID           int64  `json:"gameId" db:"match_id"`
	StartedAt    *int64 `json:"startedAt" db:"started_at"`
	Duration     *int   `json:"duration" db:"duration"`
	QueueId      *int   `json:"queueId" db:"queue_id"`
	ReplaySynced *bool  `json:"replaySynced" db:"replay_synced"`
}

type LolMatchParticipantSummary struct {
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

// LolMatchFilter provides filtering options for listing League of Legends matches
type LolMatchFilter struct {
	// Match/Summary filters
	MatchID      *int64
	StartedAtMin *int64
	StartedAtMax *int64
	ReplaySynced *bool

	// Participant filters
	PUUID      *string
	ChampionID *int
	Lane       *string
	Win        *bool
}
