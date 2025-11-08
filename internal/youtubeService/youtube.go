package youtubeService

import (
	"context"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YoutubeService interface {
	SearchLatestVideo(ctx context.Context, channelID, query string) (string, string, error)
}

type youtubeServiceImpl struct {
	service *youtube.Service
}

func NewYoutubeService(ctx context.Context, apiKey string) (YoutubeService, error) {
	srv, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &youtubeServiceImpl{service: srv}, nil
}

func (s *youtubeServiceImpl) SearchLatestVideo(ctx context.Context, channelID, query string) (string, string, error) {
	call := s.service.Search.List([]string{"id", "snippet"}).
		ChannelId(channelID).
		Q(query).
		Type("video").
		MaxResults(1).
		Order("date")

	response, err := call.Do()
	if err != nil {
		return "", "", err
	}

	if len(response.Items) == 0 {
		return "", "", fmt.Errorf("no videos found")
	}

	item := response.Items[0]
	videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", item.Id.VideoId)
	return videoURL, item.Snippet.Title, nil
}
