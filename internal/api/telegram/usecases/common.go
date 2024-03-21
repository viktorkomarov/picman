package usecases

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/viktorkomarov/picman/internal/api/telegram"
)

type stateActionProxy struct {
	notifyFunc         func(telegram.FSMContext) error
	applyUserEventFunc func(telegram.FSMContext, <-chan telegram.UserEvent) telegram.StateResult
}

func NewStateAction(
	notifyFunc func(telegram.FSMContext) error,
	applyUserEventFunc func(telegram.FSMContext, <-chan telegram.UserEvent) telegram.StateResult,
) telegram.StateAction {
	return stateActionProxy{
		notifyFunc:         notifyFunc,
		applyUserEventFunc: applyUserEventFunc,
	}
}

func (s stateActionProxy) NotifyUser(ctx telegram.FSMContext) error {
	return s.notifyFunc(ctx)
}

func (s stateActionProxy) ApplyUserEvent(ctx telegram.FSMContext, event <-chan telegram.UserEvent) telegram.StateResult {
	return s.applyUserEventFunc(ctx, event)
}

func SendToUserMessage(bot *tgbotapi.BotAPI, text string) func(telegram.FSMContext) error {
	return func(ctx telegram.FSMContext) error {
		chatID, err := telegram.GetFromUseCaseContext[int64](ctx, "chatID")
		if err != nil {
			return err
		}
		_, err = bot.Send(tgbotapi.NewMessage(chatID, text))
		return err
	}
}

func EmptyNotifyFunc() func(telegram.FSMContext) error {
	return func(f telegram.FSMContext) error {
		return nil
	}
}

func EmptyAction() func(telegram.FSMContext, <-chan telegram.UserEvent) telegram.StateResult {
	return func(_ telegram.FSMContext, _ <-chan telegram.UserEvent) telegram.StateResult {
		return telegram.StateResult{}
	}
}

type FSMProvider interface {
	GetFSMByCommandType(_type telegram.FSMType) *telegram.FSM
}
