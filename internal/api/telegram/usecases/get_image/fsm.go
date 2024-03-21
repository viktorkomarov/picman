package getimage

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/viktorkomarov/picman/internal/api/telegram"
	"github.com/viktorkomarov/picman/internal/api/telegram/usecases"
	"github.com/viktorkomarov/picman/internal/domain"
)

const (
	getImageName telegram.State = "get_image_name"
	sendImage    telegram.State = "send_image"
	completed    telegram.State = "completed"
	panicState   telegram.State = "panic_state"
)

type getImageBox struct {
	repo domain.ImageRepository
	bot  *tgbotapi.BotAPI
	name string
}

func NewGetImageFSM(repo domain.ImageRepository, bot *tgbotapi.BotAPI) *telegram.FSM {
	box := &getImageBox{
		repo: repo,
		bot:  bot,
	}

	return telegram.NewFSM(
		getImageName,
		map[telegram.State]telegram.StateAction{
			getImageName: usecases.NewStateAction(
				usecases.SendToUserMessage(bot, "Укажите название файла"),
				box.setImageNameAction,
			),
			sendImage: usecases.NewStateAction(
				box.sendImage,
				usecases.EmptyAction(),
			),
			completed: usecases.NewStateAction(
				usecases.EmptyNotifyFunc(),
				usecases.EmptyAction(),
			),
			panicState: usecases.NewStateAction(
				usecases.SendToUserMessage(bot, "Упс, такого файла нет, проверьте все файлы"),
				usecases.EmptyAction(),
			),
		},
		map[telegram.State][]telegram.State{
			getImageName: {sendImage, panicState},
			sendImage:    {completed, panicState},
		},
		[]telegram.State{completed, panicState},
	)
}

func (u *getImageBox) setImageNameAction(fsmContex telegram.FSMContext, eventCh <-chan telegram.UserEvent) telegram.StateResult {
	event, ok := <-eventCh
	if !ok {
		return toPanicState(fmt.Errorf("expected to receive user event"))
	}

	u.name = event.RawMessage.Text
	return toNextState(sendImage)
}

func (u *getImageBox) sendImage(ctx telegram.FSMContext) error {
	image, err := u.repo.GetByName(u.name)
	if err != nil {
		return fmt.Errorf("failed to find image:%w", err)
	}
	chatID, err := telegram.GetFromUseCaseContext[int64](ctx, "chatID")
	if err != nil {
		return err
	}
	cfg := tgbotapi.NewPhoto(chatID, tgbotapi.FileBytes{
		Name:  image.Name,
		Bytes: image.Payload,
	})

	_, err = u.bot.Send(cfg)
	return err
}

func toPanicState(err error) telegram.StateResult {
	return telegram.NewStateResult(panicState, err.Error())
}

func toNextState(state telegram.State) telegram.StateResult {
	return telegram.NewStateResult(state, "")
}
