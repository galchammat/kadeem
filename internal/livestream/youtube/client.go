package youtube

import (
	"context"
)

type YoutubeClient struct {
	ctx context.Context
}

func NewYoutubeClient(ctx context.Context) *YoutubeClient {
	return &YoutubeClient{
		ctx: ctx,
	}
}
