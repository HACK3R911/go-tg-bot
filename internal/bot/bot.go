package bot

import (
	"context"
	"github.com/HACK3R911/go-tg-bot/internal/auth"
	"github.com/HACK3R911/go-tg-bot/internal/handlers"
	"github.com/HACK3R911/go-tg-bot/internal/youtubeService"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type Bot struct {
	api     *tgbotapi.BotAPI
	authSvc *auth.AuthService
	// limiter  *handlers.RateLimiter
	ytSvc       youtubeService.YoutubeService
	channelId   string
	searchQuery string
}

func NewBot(token string, authSvc *auth.AuthService, ytSvc youtubeService.YoutubeService, channelId, searchQuery string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		api:     api,
		authSvc: authSvc,
		// limiter:  limiter,
		ytSvc:       ytSvc,
		channelId:   channelId,
		searchQuery: searchQuery,
	}, nil
}

func (b *Bot) Run(ctx context.Context) {
	log.Printf("Бот авторизован на аккаунте %s", b.api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 180

	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down bot...")
			return
		case update := <-updates:
			if update.Message == nil {
				continue
			}
			cmd := update.Message.Command()
			switch cmd {
			case "start":
				handlers.HandleStart(&update, b.api, b.authSvc)
			case "snake":
				handlers.HandleSnake(&update, b.api, b.authSvc /*b.limiter,*/, b.ytSvc, b.channelId, b.searchQuery)
			default:
				tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда")
			}
		}
	}
}
