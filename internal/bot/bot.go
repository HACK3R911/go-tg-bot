package bot

import (
	"context"
	"github.com/HACK3R911/go-tg-bot/internal/handler"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

const (
	startCommand = "start"
	snakeCommand = "snake"
)

type Bot struct {
	api         *tgbotapi.BotAPI
	handler     *handler.Handler
	channelId   string
	searchQuery string
}

func NewBot(token string, handler *handler.Handler, channelId, searchQuery string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		api: api,
		// limiter:  limiter,
		handler:     handler,
		channelId:   channelId,
		searchQuery: searchQuery,
	}, nil
}

func (b *Bot) Run(ctx context.Context) error {
	log.Printf("Бот авторизован на аккаунте %s", b.api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 180

	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down bot...")
			return ctx.Err()
		case update := <-updates:
			if update.Message == nil {
				continue
			}
			go func(u tgbotapi.Update) {
				if u.Message != nil && u.Message.IsCommand() {
					cmd := u.Message.Command()
					switch cmd {
					case startCommand:
						b.handler.HandleStart(&update, b.api)

					case snakeCommand:
						b.handler.HandleSnake(&update, b.api, b.channelId, b.searchQuery)
					default:
						tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда")
					}
				}
			}(update)
		}
	}
}
