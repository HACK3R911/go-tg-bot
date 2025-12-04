package main

import (
	"context"
	"flag"
	"github.com/HACK3R911/go-tg-bot/internal/handler"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"log"
	"os"
	"os/signal"
	"syscall"

	youtubeclient "github.com/HACK3R911/go-tg-bot/internal/adapter/youtube"
	"github.com/HACK3R911/go-tg-bot/internal/bot"
	"github.com/HACK3R911/go-tg-bot/internal/config"
	"github.com/HACK3R911/go-tg-bot/internal/config/env"
	"github.com/HACK3R911/go-tg-bot/internal/repository"
	"github.com/HACK3R911/go-tg-bot/internal/service"
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

	googleYT, err := youtube.NewService(ctx,
		option.WithAPIKey(apiConfig.YoutubeApiKey()),
	)
	if err != nil {
		log.Fatalf("Ошибка при создании YouTube API клиента: %v", err)
	}
	ytClient := youtubeclient.NewYoutubeAdapter(googleYT)

	repo := repository.NewRepository()
	service := service.NewService(repo, ytClient)
	handler := handler.NewHandler(service)

	tgBot, err := bot.NewBot(apiConfig.TelegramBotToken(), handler, apiConfig.ChannelId(), apiConfig.SearchQuery())
	if err != nil {
		log.Fatalf("Ошибка при создании бота: %v", err)
	}

	go tgBot.Run(ctx)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
}
