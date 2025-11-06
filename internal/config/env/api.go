package env

import (
	"errors"
	"github.com/HACK3R911/go-tg-bot/internal/config"
	"os"
)

const (
	channelIdEnvName        = "CHANNEL_ID"
	searchQueryEnvName      = "SEARCH_QUERY"
	telegramBotTokenEnvName = "TELEGRAM_BOT_TOKEN"
	youtubeApiKeyEnvName    = "YOUTUBE_API_KEY"
)

var _ config.APIConfig = (*apiConfig)(nil)

type apiConfig struct {
	channelId        string
	searchQuery      string
	telegramBotToken string
	youtubeApiKey    string
}

func NewAPIConfig() (*apiConfig, error) {
	channelId := os.Getenv(channelIdEnvName)
	if len(channelId) == 0 {
		return nil, errors.New("CHANNEL_ID не найден")
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

	return &apiConfig{
		channelId:        channelId,
		searchQuery:      searchQuery,
		telegramBotToken: telegramBotToken,
		youtubeApiKey:    youtubeApiKey,
	}, nil
}

func (cfg *apiConfig) GetChannelId() string {
	return cfg.channelId
}

func (cfg *apiConfig) GetSearchQuery() string {
	return cfg.searchQuery
}

func (cfg *apiConfig) GetTelegramBotToken() string {
	return cfg.telegramBotToken
}

func (cfg *apiConfig) GetYoutubeApiKey() string {
	return cfg.youtubeApiKey
}
