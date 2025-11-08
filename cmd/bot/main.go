package main

import (
	"context"
	"flag"
	"github.com/HACK3R911/go-tg-bot/internal/auth"
	"github.com/HACK3R911/go-tg-bot/internal/bot"
	"github.com/HACK3R911/go-tg-bot/internal/config"
	"github.com/HACK3R911/go-tg-bot/internal/config/env"
	"github.com/HACK3R911/go-tg-bot/internal/youtubeService"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "путь к конфигурационному файлу")
}

func main() {
	flag.Parse()
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Ошибка при загрузке конфигурации: %v", err)
	}

	apiConfig, err := env.NewAPIConfig()
	if err != nil {
		log.Fatalf("Ошибка при создании API конфигурации: %v", err)
	}

	ytSvc, err := youtubeService.NewYoutubeService(ctx, apiConfig.YoutubeApiKey())
	if err != nil {
		log.Fatalf("Ошибка при создании YouTube сервиса: %v", err)
	}

	authSvc := auth.NewAuthService()

	tgBot, err := bot.NewBot(apiConfig.TelegramBotToken(), authSvc, ytSvc, apiConfig.ChannelId(), apiConfig.SearchQuery())
	if err != nil {
		log.Fatalf("Ошибка при создании бота: %v", err)
	}

	go tgBot.Run(ctx)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
}
