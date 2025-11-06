package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/HACK3R911/go-tg-bot/internal/config"
	"github.com/HACK3R911/go-tg-bot/internal/config/env"
	"log"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

//const (
//	CHANNEL_ID         = ""
//	SEARCH_QUERY       = "змей"
//	TELEGRAM_BOT_TOKEN = ""
//	YOUTUBE_API_KEY    = ""
//)

var configPath string

type UserLastCommand struct {
	LastUsage time.Time
	Count     int
}

var (
	userCommands = make(map[int64]*UserLastCommand)
	commandMutex sync.RWMutex

	// Настройки ограничений
	rateLimit = 3           // Максимальное количество команд
	cooldown  = time.Minute // Период сброса ограничения

	// Добавляем глобальную мапу для хранения авторизованных пользователей
	authorizedUsers = make(map[int64]bool)
	authMutex       sync.RWMutex
)

func isCommandAllowed(userID int64) bool {
	commandMutex.Lock()
	defer commandMutex.Unlock()

	now := time.Now()
	if cmd, exists := userCommands[userID]; exists {
		// Сброс счетчика если прошло время кулдауна
		if now.Sub(cmd.LastUsage) > cooldown {
			cmd.Count = 1
			cmd.LastUsage = now
			return true
		}

		// Проверка лимита
		if cmd.Count >= rateLimit {
			return false
		}

		cmd.Count++
		cmd.LastUsage = now
		return true
	}

	// Первое использование команды этим пользователем
	userCommands[userID] = &UserLastCommand{
		LastUsage: now,
		Count:     1,
	}
	return true
}

// Функция проверки авторизации пользователя
func isUserAuthorized(userID int64) bool {
	authMutex.RLock()
	defer authMutex.RUnlock()
	return authorizedUsers[userID]
}

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

	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(apiConfig.GetYoutubeApiKey()))
	if err != nil {
		log.Fatalf("Ошибка при создании YouTube клиента: %v", err)
	}

	// Инициализация Telegram бота
	bot, err := tgbotapi.NewBotAPI(apiConfig.GetTelegramBotToken())
	if err != nil {
		log.Fatalf("Ошибка при создании Telegram бота: %v", err)
	} else {
		log.Printf("Авторизация Telegram бота успешно выполнена")
	}

	// Настройка обновлений
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 180
	updateConfig.Offset = -1

	updates := bot.GetUpdatesChan(updateConfig)

	// Обработка сообщений
	for update := range updates {
		if update.Message == nil {
			continue
		}

		userID := update.Message.From.ID

		// Обработка команды /start
		if update.Message.Command() == "start" && update.Message.Chat.Type == "private" {
			authMutex.Lock()
			authorizedUsers[userID] = true
			authMutex.Unlock()

			msg := tgbotapi.NewMessage(userID, "Вы успешно авторизованы! Теперь вы можете использовать команду /snake в любом чате.")
			bot.Send(msg)
			continue
		}

		// Проверка авторизации для команды snake
		if update.Message.Command() == "snake" {
			if !isUserAuthorized(userID) {
				if !isCommandAllowed(userID) {
					message := "Пожалуйста, подождите 3 минуты перед следующим использованием команды"

					// Пробуем отправить в личку
					privateMsg := tgbotapi.NewMessage(userID, message)
					if _, err := bot.Send(privateMsg); err != nil {
						if _, err := bot.Send(privateMsg); err != nil {
							log.Printf("Ошибка при отправке сообщения: %v", err)
						}
					}
					continue
				}

				// Поиск последнего видео со змеем
				call := youtubeService.Search.List([]string{"id", "snippet"}).
					ChannelId(apiConfig.GetChannelId()).
					Q(apiConfig.GetSearchQuery()).
					Type("video").
					MaxResults(1).
					Order("date") // Сортировка по дате (самое новое первым)

				response, err := call.Do()
				if err != nil {
					log.Printf("Ошибка при поиске видео: %v", err)
					continue
				}

				if len(response.Items) > 0 {
					item := response.Items[0]
					videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", item.Id.VideoId)
					message := fmt.Sprintf("Последнее видео со змеем:\n%s\nНазвание: %s",
						videoURL,
						item.Snippet.Title,
					)

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
					if _, err := bot.Send(msg); err != nil {
						log.Printf("Ошибка при отправке сообщения: %v", err)
					}
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Видео со змеем не найдены")
					bot.Send(msg)
				}
			}
		}
	}
}
