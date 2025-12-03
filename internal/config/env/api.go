package env

import (
	"errors"
	"github.com/HACK3R911/go-tg-bot/internal/config"
	"os"
)

const (
	ytChannelIdEnvName      = "YT_CHANNEL_ID"
	searchQueryEnvName      = "SEARCH_QUERY"
	telegramBotTokenEnvName = "TELEGRAM_BOT_TOKEN"
	youtubeApiKeyEnvName    = "YOUTUBE_API_KEY"
)

var _ config.APIConfig = (*apiConfig)(nil)

type apiConfig struct {
	ytChannelId      string
	searchQuery      string
	telegramBotToken string
	youtubeApiKey    string
	dsn              string
}

func NewAPIConfig() (*apiConfig, error) {
	ytChannelId := os.Getenv(ytChannelIdEnvName)
	if len(ytChannelId) == 0 {
		return nil, errors.New("YT_CHANNEL_ID не найден")
	}

	searchQuery := os.Getenv(searchQueryEnvName)
	if len(searchQuery) == 0 {
		return nil, errors.New("SEARCH_QUERY не найден")
	}

	telegramBotToken := os.Getenv(telegramBotTokenEnvName)
	if len(telegramBotToken) == 0 {
		return nil, errors.New("TELEGRAM_BOT_TOKEN не найден")
	}

	youtubeApiKey := os.Getenv(youtubeApiKeyEnvName)
	if len(youtubeApiKey) == 0 {
		return nil, errors.New("YOUTUBE_API_KEY не найден")
	}

	dsn := os.Getenv("DSN_TEST")
	if len(dsn) == 0 {
		return nil, errors.New("DSN_TEST не найден")
	}

	return &apiConfig{
		ytChannelId:      ytChannelId,
		searchQuery:      searchQuery,
		telegramBotToken: telegramBotToken,
		youtubeApiKey:    youtubeApiKey,
		dsn:              dsn,
	}, nil
}

func (cfg *apiConfig) ChannelId() string {
	return cfg.ytChannelId
}

func (cfg *apiConfig) SearchQuery() string {
	return cfg.searchQuery
}

func (cfg *apiConfig) TelegramBotToken() string {
	return cfg.telegramBotToken
}

func (cfg *apiConfig) YoutubeApiKey() string {
	return cfg.youtubeApiKey
}

func (cfg *apiConfig) DSN() string {
	return cfg.dsn
}
