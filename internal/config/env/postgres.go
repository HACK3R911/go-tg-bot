package env

import (
	"errors"
	"github.com/HACK3R911/go-tg-bot/internal/config"
	"os"
)

const (
	pgHostEnvName     = "PG_HOST"
	pgPortEnvName     = "PG_PORT"
	pgUserEnvName     = "PG_USER"
	pgPasswordEnvName = "PG_PASSWORD"
	pgNameEnvName     = "PG_NAME"
	pgSslModeEnvName  = "PG_SSLMODE"
	pgMaxConnsEnvName = "PG_MAX_CONNS"
	pgTimeoutEnvName  = "PG_TIMEOUT"
)

var _ config.PGConfig = (*pgConfig)(nil)

type pgConfig struct {
	host     string
	port     string
	user     string
	password string
	name     string
	sslMode  string
	maxConns string
	timeout  string
}

func NewPGConfig() (*pgConfig, error) {
	host := os.Getenv(pgHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("PG_HOST не найден")
	}

	port := os.Getenv(pgPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("PG_PORT не найден")
	}

	user := os.Getenv(pgUserEnvName)
	if len(user) == 0 {
		return nil, errors.New("PG_USER не найден")
	}

	password := os.Getenv(pgPasswordEnvName)
	if len(password) == 0 {
		return nil, errors.New("PG_PASSWORD не найден")
	}

	name := os.Getenv(pgNameEnvName)
	if len(name) == 0 {
		return nil, errors.New("PG_NAME не найден")
	}

	sslMode := os.Getenv(pgSslModeEnvName)
	if len(sslMode) == 0 {
		return nil, errors.New("PG_SSLMODE не найден")
	}

	maxConns := os.Getenv(pgMaxConnsEnvName)
	if len(maxConns) == 0 {
		return nil, errors.New("PG_MAX_CONNS не найден")
	}

	timeout := os.Getenv(pgTimeoutEnvName)
	if len(timeout) == 0 {
		return nil, errors.New("PG_TIMEOUT не найден")
	}

	return &pgConfig{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		name:     name,
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

func (cfg *pgConfig) Name() string {
	return cfg.name
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
