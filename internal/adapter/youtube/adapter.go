package youtube

import (
	"context"
	"fmt"
	"google.golang.org/api/youtube/v3"
)

type youtubeAdapter struct {
	svc *youtube.Service
}

func NewYoutubeAdapter(svc *youtube.Service) YoutubeClient {
	return &youtubeAdapter{svc: svc}
}

func (s *youtubeAdapter) SearchLatestVideo(ctx context.Context, channelID, query string) (*Video, error) {
	call := s.svc.Search.List([]string{"id", "snippet"}).
		ChannelId(channelID).
		Q(query).
		Type("video").
		MaxResults(1).
		Order("date").
		Context(ctx)

	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	if len(response.Items) == 0 {
		return nil, fmt.Errorf("no videos found")
	}

	item := response.Items[0]
	if item.Id == nil || item.Id.VideoId == "" {
		return nil, fmt.Errorf("youtube: неправильный ответ")
	}

	video := &Video{
		ID:    item.Id.VideoId,
		Title: item.Snippet.Title,
		URL:   fmt.Sprintf("https://www.youtube.com/watch?v=%s", item.Id.VideoId),
	}
	return video, nil
}
