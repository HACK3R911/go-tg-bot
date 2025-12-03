package service

import "github.com/HACK3R911/go-tg-bot/internal/repository"

type AuthService struct {
	repo repository.AuthRepo
}

func NewAuthService(repo repository.AuthRepo) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) Authorize(userID int64) {
	s.repo.AuthorizeRepo(userID)
}

func (s *AuthService) IsAuthorized(userID int64) bool {
	return s.repo.IsAuthorizedRepo(userID)
}
