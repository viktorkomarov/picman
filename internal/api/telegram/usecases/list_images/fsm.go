package listimages

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
	"github.com/viktorkomarov/picman/internal/api/telegram"
	"github.com/viktorkomarov/picman/internal/api/telegram/usecases"
	"github.com/viktorkomarov/picman/internal/domain"
)

const (
	sendList  telegram.State = "send_list"
	completed telegram.State = "completed"
	panic     telegram.State = "panic"
)

type listImagesBox struct {
	repo domain.ImageRepository
	bot  *tgbotapi.BotAPI
}

func NewFSM(repo domain.ImageRepository, bot *tgbotapi.BotAPI) *telegram.FSM {
	box := &listImagesBox{
		repo: repo,
		bot:  bot,
	}

	return telegram.NewFSM(
		sendList,
		map[telegram.State]telegram.StateAction{
			sendList:  usecases.NewStateAction(usecases.EmptyNotifyFunc(), box.sendListImage),
			completed: usecases.NewStateAction(usecases.EmptyNotifyFunc(), usecases.EmptyAction()),
			panic:     usecases.NewStateAction(usecases.ErrorUserNotify(bot), usecases.EmptyAction()),
		},
		map[telegram.State][]telegram.State{
			sendList: {completed, panic},
		},
		[]telegram.State{completed, panic},
	)
}

func (u *listImagesBox) sendListImage(ctx telegram.FSMContext, _ <-chan telegram.UserEvent) telegram.StateResult {
	images, err := u.repo.ListImages()
	if err != nil {
		return usecases.ErrorState(panic, err)
	}

	onlyNames := lo.Map(images, func(image domain.Image, _ int) string {
		return image.Name
	})
	formatMessage := strings.Join(onlyNames, "\n")

	chatID, err := telegram.GetFromUseCaseContext[int64](ctx, "chatID")
	if err != nil {
		return usecases.ErrorState(panic, err)
	}

	_, err = u.bot.Send(tgbotapi.NewMessage(chatID, formatMessage))
	if err != nil {
		return usecases.ErrorState(panic, err)
	}
	return usecases.OkState(completed)
}
