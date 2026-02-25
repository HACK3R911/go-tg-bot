package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockQuerier struct {
	mock.Mock
}

func (m *MockQuerier) Exec(ctx context.Context, sql string, args ...any) (int64, error) {
	callArgs := m.Called(ctx, sql, args)
	return callArgs.Get(0).(int64), callArgs.Error(1)
}

func (m *MockQuerier) QueryRow(ctx context.Context, sql string, args ...any) *MockRowScanner {
	callArgs := m.Called(ctx, sql, args)
	if callArgs.Get(0) == nil {
		return nil
	}
	return callArgs.Get(0).(*MockRowScanner)
}

type MockRowScanner struct {
	mock.Mock
}

func (m *MockRowScanner) Scan(dest ...any) error {
	args := m.Called(dest)
	return args.Error(0)
}

type PostgresAuthRepoMock struct {
	querier *MockQuerier
}

func NewPostgresAuthRepoMock(q *MockQuerier) *PostgresAuthRepoMock {
	return &PostgresAuthRepoMock{querier: q}
}

func (r *PostgresAuthRepoMock) AuthorizeRepo(userID int64) {
	query := `
		INSERT INTO users (id, authorized, authorized_at, updated_at)
		VALUES ($1, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (id) DO UPDATE SET
			authorized = true,
			authorized_at = CURRENT_TIMESTAMP,
			updated_at = CURRENT_TIMESTAMP
	`
	_, _ = r.querier.Exec(context.Background(), query, userID)
}

func (r *PostgresAuthRepoMock) IsAuthorizedRepo(userID int64) bool {
	query := `SELECT authorized FROM users WHERE id = $1`
	row := r.querier.QueryRow(context.Background(), query, userID)
	if row == nil {
		return false
	}

	var authorized bool
	err := row.Scan(&authorized)
	if err != nil {
		return false
	}
	return authorized
}

func TestPostgresAuthRepo_AuthorizeRepo(t *testing.T) {
	mockQuerier := new(MockQuerier)

	mockQuerier.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(int64(1), nil)

	repo := NewPostgresAuthRepoMock(mockQuerier)
	repo.AuthorizeRepo(123)

	mockQuerier.AssertExpectations(t)
}

func TestPostgresAuthRepo_AuthorizeRepo_Error(t *testing.T) {
	mockQuerier := new(MockQuerier)

	mockQuerier.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(int64(0), errors.New("db error"))

	repo := NewPostgresAuthRepoMock(mockQuerier)
	repo.AuthorizeRepo(123)

	mockQuerier.AssertExpectations(t)
}

func TestPostgresAuthRepo_IsAuthorizedRepo(t *testing.T) {
	mockQuerier := new(MockQuerier)
	mockRow := new(MockRowScanner)

	mockRow.On("Scan", mock.Anything).Return(nil)

	mockQuerier.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(mockRow)

	repo := NewPostgresAuthRepoMock(mockQuerier)
	repo.IsAuthorizedRepo(123)

	mockQuerier.AssertExpectations(t)
	mockRow.AssertCalled(t, "Scan", mock.Anything)
}

func TestPostgresAuthRepo_IsAuthorizedRepo_NotFound(t *testing.T) {
	mockQuerier := new(MockQuerier)
	mockRow := new(MockRowScanner)

	mockRow.On("Scan", mock.Anything).Return(errors.New("no rows"))

	mockQuerier.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(mockRow)

	repo := NewPostgresAuthRepoMock(mockQuerier)
	result := repo.IsAuthorizedRepo(999)

	assert.False(t, result)
	mockQuerier.AssertExpectations(t)
}

func TestPostgresAuthRepo_IsAuthorizedRepo_NotAuthorized(t *testing.T) {
	mockQuerier := new(MockQuerier)
	mockRow := new(MockRowScanner)

	mockRow.On("Scan", mock.Anything).Return(nil)

	mockQuerier.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(mockRow)

	repo := NewPostgresAuthRepoMock(mockQuerier)
	repo.IsAuthorizedRepo(123)

	mockQuerier.AssertExpectations(t)
	mockRow.AssertCalled(t, "Scan", mock.Anything)
}

type PostgresSnakeUsageRepoMock struct {
	querier *MockQuerier
}

func NewPostgresSnakeUsageRepoMock(q *MockQuerier) *PostgresSnakeUsageRepoMock {
	return &PostgresSnakeUsageRepoMock{querier: q}
}

func (r *PostgresSnakeUsageRepoMock) IncrementSnakeCounter(userID int64) {
	query := `
		INSERT INTO snake_usage (user_id, counter, created_at, updated_at)
		VALUES ($1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (user_id) DO UPDATE SET
			counter = snake_usage.counter + 1,
			updated_at = CURRENT_TIMESTAMP
	`
	_, _ = r.querier.Exec(context.Background(), query, userID)
}

func (r *PostgresSnakeUsageRepoMock) GetSnakeCounter(userID int64) int {
	query := `SELECT counter FROM snake_usage WHERE user_id = $1`
	row := r.querier.QueryRow(context.Background(), query, userID)
	if row == nil {
		return 0
	}

	var counter int
	if err := row.Scan(&counter); err != nil {
		return 0
	}
	return counter
}

func TestPostgresSnakeUsageRepo_IncrementSnakeCounter(t *testing.T) {
	mockQuerier := new(MockQuerier)

	mockQuerier.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(int64(1), nil)

	repo := NewPostgresSnakeUsageRepoMock(mockQuerier)
	repo.IncrementSnakeCounter(123)

	mockQuerier.AssertExpectations(t)
}

func TestPostgresSnakeUsageRepo_IncrementSnakeCounter_Error(t *testing.T) {
	mockQuerier := new(MockQuerier)

	mockQuerier.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(int64(0), errors.New("db error"))

	repo := NewPostgresSnakeUsageRepoMock(mockQuerier)
	repo.IncrementSnakeCounter(123)

	mockQuerier.AssertExpectations(t)
}

func TestPostgresSnakeUsageRepo_GetSnakeCounter(t *testing.T) {
	mockQuerier := new(MockQuerier)
	mockRow := new(MockRowScanner)

	mockRow.On("Scan", mock.Anything).Return(nil)

	mockQuerier.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(mockRow)

	repo := NewPostgresSnakeUsageRepoMock(mockQuerier)
	repo.GetSnakeCounter(123)

	mockQuerier.AssertExpectations(t)
	mockRow.AssertCalled(t, "Scan", mock.Anything)
}

func TestPostgresSnakeUsageRepo_GetSnakeCounter_NotFound(t *testing.T) {
	mockQuerier := new(MockQuerier)
	mockRow := new(MockRowScanner)

	mockRow.On("Scan", mock.Anything).Return(errors.New("no rows"))

	mockQuerier.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(mockRow)

	repo := NewPostgresSnakeUsageRepoMock(mockQuerier)
	result := repo.GetSnakeCounter(999)

	assert.Equal(t, 0, result)
	mockQuerier.AssertExpectations(t)
}

func TestPostgresSnakeUsageRepo_GetSnakeCounter_Zero(t *testing.T) {
	mockQuerier := new(MockQuerier)
	mockRow := new(MockRowScanner)

	mockRow.On("Scan", mock.Anything).Return(nil)

	mockQuerier.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(mockRow)

	repo := NewPostgresSnakeUsageRepoMock(mockQuerier)
	result := repo.GetSnakeCounter(123)

	assert.Equal(t, 0, result)
	mockQuerier.AssertExpectations(t)
}
