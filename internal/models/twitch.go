package models

type ChannelSearchResponse struct {
	Data []struct {
		ID          string `json:"id"`
		DisplayName string `json:"display_name"`
		AvatarURL   string `json:"thumbnail_url"`
	} `json:"data"`
}
