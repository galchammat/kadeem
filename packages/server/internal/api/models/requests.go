package models

// Riot Account Requests
type AddAccountRequest struct {
	Region     string `json:"region"`
	GameName   string `json:"game_name"`
	TagLine    string `json:"tag_line"`
	StreamerID int    `json:"streamer_id"`
}

type UpdateAccountRequest struct {
	Region   string `json:"region"`
	GameName string `json:"game_name"`
	TagLine  string `json:"tag_line"`
}

type ListMatchesRequest struct {
	PUUID  string `json:"puuid"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

type SyncMatchReplayRequest struct {
	URL string `json:"url"`
}

type SyncMatchSummaryRequest struct {
	FullMatchID string `json:"full_match_id"`
	Region      string `json:"region"`
}

// Livestream Requests
type AddStreamerRequest struct {
	Name string `json:"name"`
}

type AddChannelRequest struct {
	StreamerID  int    `json:"streamer_id"`
	ChannelName string `json:"channel_name"`
	ChannelID   string `json:"channel_id"`
	Platform    string `json:"platform"`
}

type ListBroadcastsRequest struct {
	ChannelID string `json:"channel_id"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
}
