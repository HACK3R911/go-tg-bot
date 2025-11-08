package handlers

import (
	"context"
	"fmt"
	"github.com/HACK3R911/go-tg-bot/internal/auth"
	"github.com/HACK3R911/go-tg-bot/internal/youtubeService"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func HandleStart(update *tgbotapi.Update, bot *tgbotapi.BotAPI, authSvc *auth.AuthService) {
	userID := update.Message.From.ID
	if update.Message.Chat.Type != "private" {
		return
	}

	authSvc.Authorize(userID)
	msg := tgbotapi.NewMessage(userID, "Вы успешно авторизованы! Теперь вы можете использовать команду /snake в любом чате с ботом.")
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки: %v", err)
	}
}

func HandleSnake(update *tgbotapi.Update, bot *tgbotapi.BotAPI, authSvc *auth.AuthService /*limiter *RateLimiter,*/, ytSvc youtubeService.YoutubeService, channelId, searchQuery string) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	if !authSvc.IsAuthorized(userID) {
		privateMsg := tgbotapi.NewMessage(userID, "Вы не авторизованы. Пожалуйста, используйте /start в личном чате с ботом.")
		if _, err := bot.Send(privateMsg); err != nil {
			log.Printf("Ошибка отправки: %v", err)
		}
		return
	}

	ctx := context.Background()
	videoURL, title, err := ytSvc.SearchLatestVideo(ctx, channelId, searchQuery)
	if err != nil {
		log.Printf("Ошибка поиска видео: %v", err)
		msg := tgbotapi.NewMessage(chatID, "Ошибка при поиске видео.")
		bot.Send(msg)
		return
	}

	message := fmt.Sprintf("Последнее видео со змеем:\n%s\nНазвание: %s", videoURL, title)
	msg := tgbotapi.NewMessage(chatID, message)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки: %v", err)
	}
}
