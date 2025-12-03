package repository

import "sync"

// в Models
type authDB struct {
	users map[int64]bool
	mu    sync.RWMutex
}

func NewAuthDB() *authDB {
	return &authDB{
		users: make(map[int64]bool),
		mu:    sync.RWMutex{},
	}
}

func (s *authDB) AuthorizeRepo(userID int64) {
	s.mu.Lock()
	s.users[userID] = true
	s.mu.Unlock()
}

func (s *authDB) IsAuthorizedRepo(userID int64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.users[userID]
}
