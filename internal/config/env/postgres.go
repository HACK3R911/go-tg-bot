package env

import (
	"errors"
	"fmt"
	"os"

	"github.com/HACK3R911/go-tg-bot/internal/config"
)

const (
	pgHostEnvName     = "POSTGRES_HOST"
	pgPortEnvName     = "POSTGRES_PORT"
	pgUserEnvName     = "POSTGRES_USER"
	pgPasswordEnvName = "POSTGRES_PASSWORD"
	pgDbEnvName       = "POSTGRES_DB"
	pgSslModeEnvName  = "POSTGRES_SSLMODE"
	pgMaxConnsEnvName = "POSTGRES_MAX_CONNS"
	pgTimeoutEnvName  = "POSTGRES_TIMEOUT"
)

var _ config.PGConfig = (*pgConfig)(nil)

type pgConfig struct {
	host     string
	port     string
	user     string
	password string
	db       string
	sslMode  string
	maxConns string
	timeout  string
}

func NewPGConfig() (*pgConfig, error) {
	host := os.Getenv(pgHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("POSTGRES_HOST не найден")
	}

	port := os.Getenv(pgPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("POSTGRES_PORT не найден")
	}

	user := os.Getenv(pgUserEnvName)
	if len(user) == 0 {
		return nil, errors.New("POSTGRES_USER не найден")
	}

	password := os.Getenv(pgPasswordEnvName)
	if len(password) == 0 {
		return nil, errors.New("POSTGRES_PASSWORD не найден")
	}

	db := os.Getenv(pgDbEnvName)
	if len(db) == 0 {
		return nil, errors.New("POSTGRES_DB не найден")
	}

	sslMode := os.Getenv(pgSslModeEnvName)
	if len(sslMode) == 0 {
		return nil, errors.New("POSTGRES_SSLMODE не найден")
	}

	maxConns := os.Getenv(pgMaxConnsEnvName)
	if len(maxConns) == 0 {
		return nil, errors.New("POSTGRES_MAX_CONNS не найден")
	}

	timeout := os.Getenv(pgTimeoutEnvName)
	if len(timeout) == 0 {
		return nil, errors.New("POSTGRES_TIMEOUT не найден")
	}

	return &pgConfig{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		db:       db,
		sslMode:  sslMode,
		maxConns: maxConns,
		timeout:  timeout,
	}, nil
}

func (cfg *pgConfig) Host() string {
	return cfg.host
}

func (cfg *pgConfig) Port() string {
	return cfg.port
}

func (cfg *pgConfig) User() string {
	return cfg.user
}

func (cfg *pgConfig) Password() string {
	return cfg.password
}

func (cfg *pgConfig) DB() string {
	return cfg.db
}

func (cfg *pgConfig) SSLMode() string {
	return cfg.sslMode
}

func (cfg *pgConfig) MaxConns() string {
	return cfg.maxConns
}

func (cfg *pgConfig) Timeout() string {
	return cfg.timeout
}

func (cfg *pgConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.user,
		cfg.password,
		cfg.host,
		cfg.port,
		cfg.db,
		cfg.sslMode,
	)
}
