package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	youtubeclient "github.com/HACK3R911/go-tg-bot/internal/adapter/youtube"
	"github.com/HACK3R911/go-tg-bot/internal/bot"
	"github.com/HACK3R911/go-tg-bot/internal/config"
	"github.com/HACK3R911/go-tg-bot/internal/config/env"
	"github.com/HACK3R911/go-tg-bot/internal/handler"
	"github.com/HACK3R911/go-tg-bot/internal/repository"
	"github.com/HACK3R911/go-tg-bot/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var configPath string
var migrationsPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "путь к конфигурационному файлу")
	flag.StringVar(&migrationsPath, "migrations-path", "./migrations", "путь к файлу миграций")
}

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := config.Load(configPath); err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	apiConfig, err := env.NewAPIConfig()
	if err != nil {
		log.Fatalf("Error creating API configuration: %v", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("Error creating PostgreSQL configuration: %v", err)
	}

	// Migration up
	db, err := sql.Open("pgx", pgConfig.DSN())
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Error closing db connection on migrate: %v", err)
		}
	}(db)

	if err := goose.Up(db, migrationsPath); err != nil {
		log.Fatalf("Error migration: %v", err)
	}

	// Initialize PostgreSQL connection pool
	dbPool, err := pgxpool.New(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer dbPool.Close()

	// Verify database connection
	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
	log.Println("Successfully connected to PostgreSQL")

	googleYT, err := youtube.NewService(ctx,
		option.WithAPIKey(apiConfig.YoutubeApiKey()),
	)
	if err != nil {
		log.Fatalf("Error creating YouTube API client: %v", err)
	}
	ytClient := youtubeclient.NewYoutubeAdapter(googleYT)

	// Use PostgreSQL repository
	repo := repository.NewPostgresRepository(dbPool)
	svc := service.NewService(repo, ytClient)
	hnd := handler.NewHandler(svc)

	tgBot, err := bot.NewBot(apiConfig.TelegramBotToken(), hnd, apiConfig.ChannelId(), apiConfig.SearchQuery())
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	go func() {
		if err := tgBot.Run(ctx); err != nil {
			log.Printf("Bot error: %v", err)
		}
	}()

	// Graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	log.Println("Received shutdown signal, initiating graceful shutdown...")

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	// Cancel main context to stop bot
	cancel()

	// Wait for graceful shutdown or timeout
	select {
	case <-shutdownCtx.Done():
		log.Println("Graceful shutdown completed")
	case <-time.After(10 * time.Second):
		log.Println("Shutdown timeout exceeded, forcing exit")
	}
}
