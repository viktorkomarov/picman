package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/viktorkomarov/picman/internal/api/telegram"
)

type MessageHandler interface {
	OnRecieveUserMessage(msg telegram.Message)
}

type MessageFilter interface {
	ShouldExclude(update tgbotapi.Update) bool
}

type RunBotConfig struct {
	Token               string
	OffsetInitial       int
	PollIntervalSeconds int
	MessageFilter       MessageFilter
	MessageHandler      MessageHandler
}

func RunBot(ctx context.Context, cfg RunBotConfig) error {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return fmt.Errorf("tgbotapi.NewBotAPI: %w", err)
	}
	updateConfig := tgbotapi.NewUpdate(cfg.OffsetInitial)
	updateConfig.Timeout = cfg.PollIntervalSeconds

	updateCh := bot.GetUpdatesChan(updateConfig)
	for {
		select {
		case update := <-updateCh:
			if !cfg.MessageFilter.ShouldExclude(update) {
				// map
				cfg.MessageHandler.OnRecieveUserMessage(telegram.Message{})
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
