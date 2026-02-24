package handler

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const chatType = "private"

func (h *Handler) HandleStart(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	userID := update.Message.From.ID
	if update.Message.Chat.Type != chatType {
		return
	}

	h.service.Authorize(userID)
	msg := tgbotapi.NewMessage(userID, "Вы успешно авторизованы! Теперь вы можете использовать команду /snake в любом чате с ботом.")
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки: %v", err)
	}
	log.Printf("Пользователь %d авторизован", userID)
}

func (h *Handler) HandleSnake(update *tgbotapi.Update, bot *tgbotapi.BotAPI, channelId, searchQuery string) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	if !h.service.IsAuthorized(userID) {
		privateMsg := tgbotapi.NewMessage(userID, "Вы не авторизованы. Пожалуйста, используйте /start в личном чате с ботом.")
		if _, err := bot.Send(privateMsg); err != nil {
			log.Printf("Ошибка отправки: %v", err)
		}
		return
	}

	ctx := context.Background()

	video, err := h.service.SearchLatestVideo(ctx, channelId, searchQuery)
	if err != nil {
		log.Printf("Ошибка поиска видео: %v", err)
		msg := tgbotapi.NewMessage(chatID, "Ошибка при поиске видео.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Ошибка отправки: %v", err)
		}
		return
	}

	h.service.IncrementSnakeCounter(userID)
	count := h.service.GetSnakeCounter(userID)

	message := fmt.Sprintf("Последнее видео со змеем:\n%s\nНазвание: %s\n\nИспользовано раз: %d", video.URL, video.Title, count)

	msg := tgbotapi.NewMessage(chatID, message)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки: %v", err)
	}
}

func (h *Handler) HandleDefault(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда")
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки: %v", err)
	}
}
