package service

import (
	"context"
	youtubeclient "github.com/HACK3R911/go-tg-bot/internal/adapter/youtube"
)

type YoutubeService struct {
	ytClient youtubeclient.YoutubeClient
}

func NewYoutubeService(yt youtubeclient.YoutubeClient) *YoutubeService {
	return &YoutubeService{ytClient: yt}
}

func (s *YoutubeService) SearchLatestVideo(ctx context.Context, channelID, query string) (*youtubeclient.Video, error) {
	video, err := s.ytClient.SearchLatestVideo(ctx, channelID, query)
	if err != nil {
		return nil, err
	}
	return video, nil
}
