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
	ChannelId() string
	SearchQuery() string
	TelegramBotToken() string
	YoutubeApiKey() string
}
