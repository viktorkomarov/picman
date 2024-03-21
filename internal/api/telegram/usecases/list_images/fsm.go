package listimages

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
	"github.com/viktorkomarov/picman/internal/api/telegram"
	"github.com/viktorkomarov/picman/internal/api/telegram/usecases"
	"github.com/viktorkomarov/picman/internal/domain"
)

const (
	sendList   telegram.State = "send_list"
	completed  telegram.State = "completed"
	panicState telegram.State = "panic_state"
)

type listImagesBox struct {
	repo domain.ImageRepository
	bot  *tgbotapi.BotAPI
}

func NewListImagesFSM(repo domain.ImageRepository, bot *tgbotapi.BotAPI) *telegram.FSM {
	box := &listImagesBox{
		repo: repo,
		bot:  bot,
	}

	return telegram.NewFSM(
		sendList,
		map[telegram.State]telegram.StateAction{
			sendList: usecases.NewStateAction(
				box.sendListImage,
				usecases.EmptyAction(),
			),
			completed: usecases.NewStateAction(
				usecases.EmptyNotifyFunc(),
				usecases.EmptyAction(),
			),
			panicState: usecases.NewStateAction(
				usecases.SendToUserMessage(bot, "Что-то пошло не так - повторите попозже"),
				usecases.EmptyAction(),
			),
		},
		map[telegram.State][]telegram.State{
			sendList: {completed, panicState},
		},
		[]telegram.State{completed, panicState},
	)
}

func (u *listImagesBox) sendListImage(ctx telegram.FSMContext) error {
	images, err := u.repo.ListImages()
	if err != nil {
		return fmt.Errorf("failed to find image:%w", err)
	}

	onlyNames := lo.Map(images, func(image domain.Image, _ int) string {
		return image.Name
	})
	formatMessage := strings.Join(onlyNames, "\n")

	chatID, err := telegram.GetFromUseCaseContext[int64](ctx, "chatID")
	if err != nil {
		return err
	}

	_, err = u.bot.Send(tgbotapi.NewMessage(chatID, formatMessage))
	return err
}

func toPanicState(err error) telegram.StateResult {
	return telegram.NewStateResult(panicState, err.Error())
}

func toNextState(state telegram.State) telegram.StateResult {
	return telegram.NewStateResult(state, "")
}
