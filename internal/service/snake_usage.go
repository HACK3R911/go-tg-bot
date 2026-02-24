package service

import "github.com/HACK3R911/go-tg-bot/internal/repository"

type SnakeUsageService struct {
	repo repository.SnakeUsageRepo
}

func NewSnakeUsageService(repo repository.SnakeUsageRepo) *SnakeUsageService {
	return &SnakeUsageService{repo: repo}
}

func (s *SnakeUsageService) IncrementSnakeCounter(userID int64) {
	s.repo.IncrementSnakeCounter(userID)
}

func (s *SnakeUsageService) GetSnakeCounter(userID int64) int {
	return s.repo.GetSnakeCounter(userID)
}
