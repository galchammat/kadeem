package models

type LolApiReplaysReponse struct {
	URLs []string `json:"matchFileURLs"`
}

type LeagueOfLegendsAccount struct {
	PUUID      string `json:"puuid" db:"puuid,primarykey"`
	StreamerID int    `json:"streamerId,omitempty" db:"streamer_id"`
	TagLine    string `json:"tagLine" db:"tag_line"`
	GameName   string `json:"gameName" db:"game_name"`
	Region     string `json:"region,omitempty" db:"region"`
	SyncedAt   *int64 `json:"syncedAt" db:"synced_at"`
}

type LolAccountPatch struct {
	PUUID    string  `json:"puuid" db:"puuid,primarykey"`
	TagLine  *string `json:"tagLine,omitempty" db:"tag_line"`
	GameName *string `json:"gameName,omitempty" db:"game_name"`
	Region   *string `json:"region,omitempty" db:"region"`
	SyncedAt *int64  `json:"syncedAt,omitempty" db:"synced_at"`
}

type LeagueOfLegendsMatch struct {
	Summary      LeagueOfLegendsMatchSummary              `json:"summary" db:"-"`
	Participants []LeagueOfLegendsMatchParticipantSummary `json:"participants" db:"-"`
	ReplayURL    *string                                  `json:"replay,omitempty" db:"replay"`
}

type LeagueOfLegendsMatchSummary struct {
	ID           int64  `json:"gameId" db:"match_id"`
	StartedAt    *int64 `json:"startedAt" db:"started_at"`
	Duration     *int   `json:"duration" db:"duration"`
	ReplaySynced *bool  `json:"replaySynced" db:"replay_synced"`
}

type LeagueOfLegendsMatchParticipantSummary struct {
	GameID                      int64  `json:"gameId" db:"match_id"`
	ChampionID                  int    `json:"championId" db:"champion_id"`
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
	MatchID      *int64 `db:"m.match_id" op:"="`
	StartedAtMin *int64 `db:"m.started_at" op:">="`
	StartedAtMax *int64 `db:"m.started_at" op:"<="`
	DurationMin  *int   `db:"m.duration" op:">="`
	DurationMax  *int   `db:"m.duration" op:"<="`
	ReplaySynced *bool  `db:"m.replay_synced" op:"="`

	// Participant filters
	PUUID                          *string `db:"p.puuid" op:"="`
	Region                         *string `db:"p.region" op:"="`
	ChampionID                     *int    `db:"p.champion_id" op:"="`
	Lane                           *string `db:"p.lane" op:"="`
	Win                            *bool   `db:"p.win" op:"="`
	ParticipantID                  *int    `db:"p.participant_id" op:"="`
	RiotIDGameName                 *string `db:"p.riot_id_game_name" op:"="`
	RiotIDTagline                  *string `db:"p.riot_id_tagline" op:"="`
	MinKills                       *int    `db:"p.kills" op:">="`
	MaxKills                       *int    `db:"p.kills" op:"<="`
	MinDeaths                      *int    `db:"p.deaths" op:">="`
	MaxDeaths                      *int    `db:"p.deaths" op:"<="`
	MinAssists                     *int    `db:"p.assists" op:">="`
	MaxAssists                     *int    `db:"p.assists" op:"<="`
	MinDoubleKills                 *int    `db:"p.double_kills" op:">="`
	MaxDoubleKills                 *int    `db:"p.double_kills" op:"<="`
	MinTripleKills                 *int    `db:"p.triple_kills" op:">="`
	MaxTripleKills                 *int    `db:"p.triple_kills" op:"<="`
	MinQuadraKills                 *int    `db:"p.quadra_kills" op:">="`
	MaxQuadraKills                 *int    `db:"p.quadra_kills" op:"<="`
	MinPentaKills                  *int    `db:"p.penta_kills" op:">="`
	MaxPentaKills                  *int    `db:"p.penta_kills" op:"<="`
	MinTotalDamageDealtToChampions *int    `db:"p.total_damage_dealt_to_champions" op:">="`
	MaxTotalDamageDealtToChampions *int    `db:"p.total_damage_dealt_to_champions" op:"<="`
	MinTotalDamageTaken            *int    `db:"p.total_damage_taken" op:">="`
	MaxTotalDamageTaken            *int    `db:"p.total_damage_taken" op:"<="`
	MinTotalMinionsKilled          *int    `db:"p.total_minions_killed" op:">="`
	MaxTotalMinionsKilled          *int    `db:"p.total_minions_killed" op:"<="`
	Item0                          *int    `db:"p.item0" op:"="`
	Item1                          *int    `db:"p.item1" op:"="`
	Item2                          *int    `db:"p.item2" op:"="`
	Item3                          *int    `db:"p.item3" op:"="`
	Item4                          *int    `db:"p.item4" op:"="`
	Item5                          *int    `db:"p.item5" op:"="`
	Item6                          *int    `db:"p.item6" op:"="`
	Summoner1ID                    *int    `db:"p.summoner1_id" op:"="`
	Summoner2ID                    *int    `db:"p.summoner2_id" op:"="`
}
