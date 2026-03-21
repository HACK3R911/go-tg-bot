package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

func Load(path string) error {
	if path == "" || path == ".env" {
		if _, err := os.Stat(".env"); os.IsNotExist(err) {
			if os.Getenv("TELEGRAM_BOT_TOKEN") != "" {
				log.Println("Using environment variables (no .env file found)")
				return nil
			}
		}
	}

	if path != "" && path != ".env" || fileExists(path) {
		err := godotenv.Load(path)
		if err != nil {
			log.Printf("Ошибка загрузки .env файла: %v", err)
			return err
		}
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

type APIConfig interface {
	ChannelId() string
	SearchQuery() string
	TelegramBotToken() string
	YoutubeApiKey() string
}

type PGConfig interface {
	Host() string
	Port() string
	User() string
	Password() string
	DB() string
	SSLMode() string
	MaxConns() string
	Timeout() string
	DSN() string
}

type RLConfig interface {
	MaxConcurrent() int
	GlobalRPS() rate.Limit
	GlobalBurst() int
	PerUserRPS() rate.Limit
	PerUserBurst() int
	MaxUsers() int
}
