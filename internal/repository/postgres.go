package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresAuthRepo implements AuthRepo interface using PostgreSQL
type PostgresAuthRepo struct {
	pool *pgxpool.Pool
}

// NewPostgresAuthRepo creates a new PostgreSQL auth repository
func NewPostgresAuthRepo(pool *pgxpool.Pool) *PostgresAuthRepo {
	return &PostgresAuthRepo{pool: pool}
}

// AuthorizeRepo authorizes a user in the database
func (r *PostgresAuthRepo) AuthorizeRepo(userID int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO users (id, authorized, authorized_at, updated_at)
		VALUES ($1, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (id) DO UPDATE SET
			authorized = true,
			authorized_at = CURRENT_TIMESTAMP,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		// Log error but don't fail - auth failure shouldn't crash the bot
		fmt.Printf("Error authorizing user %d: %v\n", userID, err)
	}
}

// IsAuthorizedRepo checks if a user is authorized in the database
func (r *PostgresAuthRepo) IsAuthorizedRepo(userID int64) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT authorized FROM users WHERE id = $1`

	var authorized bool
	err := r.pool.QueryRow(ctx, query, userID).Scan(&authorized)
	if err != nil {
		// User not found or error means not authorized
		return false
	}

	return authorized
}

// GetUserStats returns statistics for a user
func (r *PostgresAuthRepo) GetUserStats(userID int64) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT authorized, authorized_at, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var authorized bool
	var authorizedAt, createdAt, updatedAt time.Time

	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&authorized,
		&authorizedAt,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting user stats: %w", err)
	}

	return map[string]interface{}{
		"authorized":    authorized,
		"authorized_at": authorizedAt,
		"created_at":    createdAt,
		"updated_at":    updatedAt,
	}, nil
}

type PostgresSnakeUsageRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresSnakeUsageRepo(pool *pgxpool.Pool) *PostgresSnakeUsageRepo {
	return &PostgresSnakeUsageRepo{pool: pool}
}

func (r *PostgresSnakeUsageRepo) IncrementSnakeCounter(userID int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO snake_usage (user_id, counter, created_at, updated_at)
		VALUES ($1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (user_id) DO UPDATE SET
			counter = snake_usage.counter + 1,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		fmt.Printf("Error incrementing snake counter for user %d: %v\n", userID, err)
	}
}

func (r *PostgresSnakeUsageRepo) GetSnakeCounter(userID int64) int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT counter FROM snake_usage WHERE user_id = $1`

	var counter int
	err := r.pool.QueryRow(ctx, query, userID).Scan(&counter)
	if err != nil {
		return 0
	}

	return counter
}
