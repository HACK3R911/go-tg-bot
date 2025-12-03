package main

import (
	"context"
	"flag"
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

//var goosePath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "путь к конфигурационному файлу")
	//flag.StringVar(&goosePath, "dir", "./migrations", "путь к файлу миграции")
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

	//pgConfig, err := env.NewPGConfig()
	//if err != nil {
	//	log.Fatalf("Ошибка при создании POSTGRES конфигурации: %v", err)
	//}

	//db, err := goose.OpenDBWithDriver("postgres", apiConfig.DSN())
	//if err != nil {
	//	log.Fatalf("goose: failed to open DB: %v", err)
	//}
	//
	//defer func() {
	//	if err := db.Close(); err != nil {
	//		log.Fatalf("goose: failed to close DB: %v", err)
	//	}
	//}()
	//
	//if err := goose.RunContext(ctx, command, db, *dir, arguments...); err != nil {
	//	log.Fatalf("goose %v: %v", command, err)
	//}

	//ytSvc, err := service.NewYoutubeService(ctx, apiConfig.YoutubeApiKey())
	//if err != nil {
	//	log.Fatalf("Ошибка при создании YouTube сервиса: %v", err)
	//}
	//authSvc := service.NewAuthService()

	googleYT, err := youtube.NewService(ctx,
		option.WithAPIKey(apiConfig.YoutubeApiKey()),
	)
	if err != nil {
		log.Fatalf("Ошибка при создании YouTube API клиента: %v", err)
	}
	ytClient := youtubeclient.NewYoutubeAdapter(googleYT)

	repo := repository.NewRepository()
	service := service.NewService(repo, ytClient)

	tgBot, err := bot.NewBot(apiConfig.TelegramBotToken(), service, apiConfig.ChannelId(), apiConfig.SearchQuery())
	if err != nil {
		log.Fatalf("Ошибка при создании бота: %v", err)
	}

	go tgBot.Run(ctx)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
}
