package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/viktorkomarov/picman/internal/api/telegram"
)

type MessageHandler interface {
	OnRecieveUserMessage(msg telegram.UserEvent)
}

type RunBotConfig struct {
	MessageHandler MessageHandler
}

func RunBot(
	ctx context.Context,
	cfg tgbotapi.UpdateConfig,
	bot *tgbotapi.BotAPI,
	receiver MessageHandler,
) error {
	updateCh := bot.GetUpdatesChan(cfg)
	for {
		select {
		case update := <-updateCh:
			if update.Message == nil {
				continue
			}
			receiver.OnRecieveUserMessage(telegram.UserEvent{
				RawMessage: update.Message,
			})
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
