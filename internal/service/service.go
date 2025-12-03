package service

import (
	"context"
	youtubeclient "github.com/HACK3R911/go-tg-bot/internal/adapter/youtube"
	"github.com/HACK3R911/go-tg-bot/internal/repository"
)

type Auth interface {
	Authorize(userID int64)
	IsAuthorized(userID int64) bool
}

type Youtube interface {
	SearchLatestVideo(ctx context.Context, channelID, query string) (*youtubeclient.Video, error)
}

type Service struct {
	Auth
	Youtube
}

func NewService(repos *repository.Repository, ytClient youtubeclient.YoutubeClient) *Service {
	return &Service{
		Auth:    NewAuthService(repos.AuthRepo),
		Youtube: NewYoutubeService(ytClient),
	}
}
