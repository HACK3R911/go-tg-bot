package env

import (
	"errors"
	"os"
	"strconv"

	"github.com/HACK3R911/go-tg-bot/internal/config"
	"golang.org/x/time/rate"
)

const (
	mlMaxConcurrentEnvName = "RATE_LIMIT_MAX_CONCURRENT"
	rlGlobalRPSEnvName     = "RATE_LIMIT_GLOBAL_RPS"
	rlGlobalBurstEnvName   = "RATE_LIMIT_GLOBAL_BURST"
	rlPerUserRPSEnvName    = "RATE_LIMIT_PER_USER_RPS"
	rlPerUserBurstEnvName  = "RATE_LIMIT_PER_USER_BURST"
	rlMaxUsersEnvName      = "RATE_LIMIT_MAX_USERS"
)

var _ config.RLConfig = (*rlConfig)(nil)

type rlConfig struct {
	maxConcurrent int
	globalRPS     rate.Limit
	globalBurst   int
	perUserRPS    rate.Limit
	perUserBurst  int
	maxUsers      int
}

func NewRLConfig() (*rlConfig, error) {
	cfg := &rlConfig{
		maxConcurrent: 10,
		globalRPS:     rate.Limit(30),
		globalBurst:   30,
		perUserRPS:    rate.Limit(3),
		perUserBurst:  5,
		maxUsers:      10000,
	}

	maxConcurrent := os.Getenv(mlMaxConcurrentEnvName)
	if maxConcurrent != "" {
		n, err := strconv.Atoi(maxConcurrent)
		if err != nil {
			return nil, errors.New("RATE_LIMIT_MAX_CONCURRENT must be a number")
		}
		cfg.maxConcurrent = n
	}

	if val := os.Getenv(rlGlobalRPSEnvName); val != "" {
		n, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return nil, errors.New("RATE_LIMIT_GLOBAL_RPS must be a number")
		}
		cfg.globalRPS = rate.Limit(n)
	}

	if val := os.Getenv(rlGlobalBurstEnvName); val != "" {
		n, err := strconv.Atoi(val)
		if err != nil {
			return nil, errors.New("RATE_LIMIT_GLOBAL_BURST must be a number")
		}
		cfg.globalBurst = n
	}

	if val := os.Getenv(rlPerUserRPSEnvName); val != "" {
		n, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return nil, errors.New("RATE_LIMIT_PER_USER_RPS must be a number")
		}
		cfg.perUserRPS = rate.Limit(n)
	}

	if val := os.Getenv(rlPerUserBurstEnvName); val != "" {
		n, err := strconv.Atoi(val)
		if err != nil {
			return nil, errors.New("RATE_LIMIT_PER_USER_BURST must be a number")
		}
		cfg.perUserBurst = n
	}

	if val := os.Getenv(rlMaxUsersEnvName); val != "" {
		n, err := strconv.Atoi(val)
		if err != nil {
			return nil, errors.New("RATE_LIMIT_MAX_USERS must be a number")
		}
		cfg.maxUsers = n
	}

	return cfg, nil
}

func (cfg *rlConfig) MaxConcurrent() int {
	return cfg.maxConcurrent
}

func (cfg *rlConfig) GlobalRPS() rate.Limit {
	return cfg.globalRPS
}

func (cfg *rlConfig) GlobalBurst() int {
	return cfg.globalBurst
}

func (cfg *rlConfig) PerUserRPS() rate.Limit {
	return cfg.perUserRPS
}

func (cfg *rlConfig) PerUserBurst() int {
	return cfg.perUserBurst
}

func (cfg *rlConfig) MaxUsers() int {
	return cfg.maxUsers
}
