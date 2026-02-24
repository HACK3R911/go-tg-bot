package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// AuthRepo defines the interface for user authorization storage
type AuthRepo interface {
	AuthorizeRepo(userID int64)
	IsAuthorizedRepo(userID int64) bool
}

// SnakeUsageRepo defines the interface for snake command usage tracking
type SnakeUsageRepo interface {
	IncrementSnakeCounter(userID int64)
	GetSnakeCounter(userID int64) int
}

// Repository aggregates all repository interfaces
type Repository struct {
	AuthRepo
	SnakeUsageRepo
}

// NewRepository creates a new repository with in-memory storage (for testing)
func NewRepository() *Repository {
	return &Repository{
		AuthRepo:       NewAuthDB(),
		SnakeUsageRepo: NewSnakeUsageDB(),
	}
}

// NewPostgresRepository creates a new repository with PostgreSQL storage
func NewPostgresRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		AuthRepo:       NewPostgresAuthRepo(pool),
		SnakeUsageRepo: NewPostgresSnakeUsageRepo(pool),
	}
}
