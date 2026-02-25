//go:build integration

package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func setupTestDB(t *testing.T) (*pgxpool.Pool, func()) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(&testcontainers.WaitStrategyStrategy{
			WaitForStrategy: testcontainers.WaitForLog("database system is ready to accept connections"),
		}),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	dsn := pgContainer.ConnString()
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("failed to ping: %v", err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id BIGINT PRIMARY KEY,
		authorized BOOLEAN DEFAULT false,
		authorized_at TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS snake_usage (
		user_id BIGINT PRIMARY KEY,
		counter INTEGER DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := pool.Exec(ctx, schema); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	cleanup := func() {
		pool.Close()
		pgContainer.Terminate(ctx)
	}

	return pool, cleanup
}

func BenchmarkPostgresAuthRepo_AuthorizeRepo(b *testing.B) {
	pool, cleanup := setupTestDB(b)
	defer cleanup()

	repo := NewPostgresAuthRepo(pool)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.AuthorizeRepo(int64(i))
	}

	_ = ctx
}

func BenchmarkPostgresAuthRepo_IsAuthorizedRepo(b *testing.B) {
	pool, cleanup := setupTestDB(b)
	defer cleanup()

	repo := NewPostgresAuthRepo(pool)

	for i := 0; i < 100; i++ {
		repo.AuthorizeRepo(int64(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.IsAuthorizedRepo(int64(i % 100))
	}
}

func BenchmarkPostgresSnakeUsageRepo_IncrementSnakeCounter(b *testing.B) {
	pool, cleanup := setupTestDB(b)
	defer cleanup()

	repo := NewPostgresSnakeUsageRepo(pool)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.IncrementSnakeCounter(int64(i % 100))
	}
}

func BenchmarkPostgresSnakeUsageRepo_GetSnakeCounter(b *testing.B) {
	pool, cleanup := setupTestDB(b)
	defer cleanup()

	repo := NewPostgresSnakeUsageRepo(pool)

	for i := 0; i < 100; i++ {
		repo.IncrementSnakeCounter(int64(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.GetSnakeCounter(int64(i % 100))
	}
}

func TestPostgresAuthRepo_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewPostgresAuthRepo(pool)
	ctx := context.Background()

	repo.AuthorizeRepo(123)

	authorized := repo.IsAuthorizedRepo(123)
	if !authorized {
		t.Error("expected user to be authorized")
	}

	notAuthorized := repo.IsAuthorizedRepo(999)
	if notAuthorized {
		t.Error("expected user to not be authorized")
	}

	stats, err := repo.GetUserStats(123)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats["authorized"] != true {
		t.Errorf("expected authorized=true, got %v", stats["authorized"])
	}

	_ = ctx
}

func TestPostgresSnakeUsageRepo_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewPostgresSnakeUsageRepo(pool)
	ctx := context.Background()

	repo.IncrementSnakeCounter(123)
	repo.IncrementSnakeCounter(123)
	repo.IncrementSnakeCounter(123)

	count := repo.GetSnakeCounter(123)
	if count != 3 {
		t.Errorf("expected counter=3, got %d", count)
	}

	notFound := repo.GetSnakeCounter(999)
	if notFound != 0 {
		t.Errorf("expected counter=0 for non-existent user, got %d", notFound)
	}

	_ = ctx
}

func ExamplePostgresAuthRepo() {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(&testcontainers.WaitStrategyStrategy{
			WaitForStrategy: testcontainers.WaitForLog("database system is ready to accept connections"),
		}),
	)
	if err != nil {
		fmt.Printf("failed to start container: %v\n", err)
		return
	}
	defer pgContainer.Terminate(ctx)

	dsn := pgContainer.ConnString()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		fmt.Printf("failed to create pool: %v\n", err)
		return
	}
	defer pool.Close()

	repo := NewPostgresAuthRepo(pool)
	repo.AuthorizeRepo(123)

	fmt.Println(repo.IsAuthorizedRepo(123))
}
