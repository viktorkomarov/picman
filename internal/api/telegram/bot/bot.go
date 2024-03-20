package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/viktorkomarov/picman/internal/api/telegram"
)

type MessageHandler interface {
	OnRecieveUserMessage(msg telegram.UserEvent)
}

type MessageFilter interface {
	ShouldExclude(update tgbotapi.Update) bool
}

type RunBotConfig struct {
	OffsetInitial       int
	PollIntervalSeconds int
	MessageFilter       MessageFilter
	MessageHandler      MessageHandler
}

func RunBot(ctx context.Context, bot *tgbotapi.BotAPI, cfg RunBotConfig) error {
	updateConfig := tgbotapi.NewUpdate(cfg.OffsetInitial)
	updateConfig.Timeout = cfg.PollIntervalSeconds

	updateCh := bot.GetUpdatesChan(updateConfig)
	for {
		select {
		case update := <-updateCh:
			cfg.MessageHandler.OnRecieveUserMessage(telegram.UserEvent{
				RawMessage: update.Message,
			})
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
