package config

import (
	"github.com/joho/godotenv"
	"log"
)

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
	return nil
}

type APIConfig interface {
	GetChannelId() string
	GetSearchQuery() string
	GetTelegramBotToken() string
	GetYoutubeApiKey() string
}
