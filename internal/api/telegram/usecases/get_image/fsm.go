package getimage

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/viktorkomarov/picman/internal/api/telegram"
	"github.com/viktorkomarov/picman/internal/api/telegram/usecases"
	"github.com/viktorkomarov/picman/internal/domain"
)

const (
	setImageName telegram.State = "set_image_name"
	sendImage    telegram.State = "send_image"
	completed    telegram.State = "completed"
	panicState   telegram.State = "panic_state"
)

type getImageBox struct {
	repo domain.ImageRepository
	bot  *tgbotapi.BotAPI
	name string
}

func NewFSM(repo domain.ImageRepository, bot *tgbotapi.BotAPI) *telegram.FSM {
	box := &getImageBox{
		repo: repo,
		bot:  bot,
	}

	return telegram.NewFSM(
		setImageName,
		map[telegram.State]telegram.StateAction{
			setImageName: usecases.NewStateAction(
				usecases.SendMessageNotifyFunc(bot, "Укажите название файла"),
				usecases.ActionWithEvent(panicState, box.setImageNameAction),
			),
			sendImage:  usecases.NewStateAction(usecases.EmptyNotifyFunc(), box.sendImage),
			completed:  usecases.NewStateAction(usecases.EmptyNotifyFunc(), usecases.EmptyAction()),
			panicState: usecases.NewStateAction(usecases.ErrorUserNotify(bot), usecases.EmptyAction()),
		},
		map[telegram.State][]telegram.State{
			setImageName: {sendImage, panicState},
			sendImage:    {completed, panicState},
		},
		[]telegram.State{completed, panicState},
	)
}

func (u *getImageBox) setImageNameAction(fsmContex telegram.FSMContext, event telegram.UserEvent) telegram.StateResult {
	u.name = event.RawMessage.Text
	return usecases.OkState(sendImage)
}

func (u *getImageBox) sendImage(ctx telegram.FSMContext, _ <-chan telegram.UserEvent) telegram.StateResult {
	image, err := u.repo.GetByName(u.name)
	if err != nil {
		return usecases.ErrorState(panicState, err)
	}
	chatID, err := telegram.GetFromUseCaseContext[int64](ctx, "chatID")
	if err != nil {
		return usecases.ErrorState(panicState, err)
	}
	cfg := tgbotapi.NewPhoto(chatID, tgbotapi.FileBytes{
		Name:  image.Name,
		Bytes: image.Payload,
	})

	if _, err = u.bot.Send(cfg); err != nil {
		return usecases.ErrorState(panicState, err)
	}
	return usecases.OkState(completed)
}
