package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type UserEvent struct {
	RawMessage *tgbotapi.Message
}

func (m UserEvent) IsCommand() bool {
	return m.RawMessage.IsCommand()
}
