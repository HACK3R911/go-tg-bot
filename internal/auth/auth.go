package auth

import "sync"

type AuthService struct {
	users map[int64]bool
	mu    sync.RWMutex
}

func NewAuthService() *AuthService {
	return &AuthService{
		users: make(map[int64]bool),
	}
}

func (s *AuthService) Authorize(userID int64) {
	s.mu.Lock()
	s.users[userID] = true
	s.mu.Unlock()
}

func (s *AuthService) IsAuthorized(userID int64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.users[userID]
}
