package youtube

import "context"

type YoutubeClient interface {
	SearchLatestVideo(ctx context.Context, channelID, query string) (*Video, error)
}

type Video struct {
	ID    string
	Title string
	URL   string
}
