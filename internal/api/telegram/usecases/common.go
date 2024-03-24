package usecases

import (
	"fmt"

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

func SendMessageNotifyFunc(bot *tgbotapi.BotAPI, text string) func(telegram.FSMContext) error {
	return func(ctx telegram.FSMContext) error {
		return SendMessage(ctx, bot, text)
	}
}

func SendMessage(ctx telegram.FSMContext, bot *tgbotapi.BotAPI, text string) error {
	chatID, err := telegram.GetFromUseCaseContext[int64](ctx, "chatID")
	if err != nil {
		return err
	}
	_, err = bot.Send(tgbotapi.NewMessage(chatID, text))
	return err
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

func ActionWithEvent(
	panicState telegram.State,
	userFunc func(telegram.FSMContext, telegram.UserEvent) telegram.StateResult,
) func(telegram.FSMContext, <-chan telegram.UserEvent) telegram.StateResult {
	return func(f telegram.FSMContext, c <-chan telegram.UserEvent) telegram.StateResult {
		event, ok := <-c
		if !ok {
			return ErrorState(panicState, fmt.Errorf("expected to receive user event"))
		}
		return userFunc(f, event)
	}
}

func ErrorUserNotify(bot *tgbotapi.BotAPI) func(telegram.FSMContext) error {
	return func(ctx telegram.FSMContext) error {
		state := ctx.LastState()
		if state.Error == nil {
			return nil
		}

		return SendMessage(ctx, bot, fmt.Sprintf("В процессе операции возникли проблемы: %s", state.Error.Error()))
	}
}

func ErrorState(next telegram.State, err error) telegram.StateResult {
	return telegram.NewStateResult(next, err)
}

func OkState(state telegram.State) telegram.StateResult {
	return telegram.NewStateResult(state, nil)
}
