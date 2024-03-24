package deleteimage

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/viktorkomarov/picman/internal/api/telegram"
	"github.com/viktorkomarov/picman/internal/api/telegram/usecases"
	"github.com/viktorkomarov/picman/internal/domain"
)

const (
	setImageName telegram.State = "set_image_name"
	deleteImage  telegram.State = "delete_image"
	completed    telegram.State = "completed"
	panic        telegram.State = "panic"
)

type deleteImageBox struct {
	repo domain.ImageRepository
	bot  *tgbotapi.BotAPI
	name string
}

func NewDeleteImageFSM(repo domain.ImageRepository, bot *tgbotapi.BotAPI) *telegram.FSM {
	box := &deleteImageBox{
		repo: repo,
		bot:  bot,
	}

	return telegram.NewFSM(
		setImageName,
		map[telegram.State]telegram.StateAction{
			setImageName: usecases.NewStateAction(
				usecases.SendMessageNotifyFunc(bot, "Укажите название файла"),
				usecases.ActionWithEvent(panic, box.setImageNameAction),
			),
			deleteImage: usecases.NewStateAction(usecases.EmptyNotifyFunc(), box.deleteImage),
			completed: usecases.NewStateAction(
				usecases.SendMessageNotifyFunc(bot, "Файл успешно удален"), usecases.EmptyAction(),
			),
			panic: usecases.NewStateAction(
				usecases.SendMessageNotifyFunc(bot, "Упс, такого файла нет, проверьте все файлы"),
				usecases.EmptyAction(),
			),
		},
		map[telegram.State][]telegram.State{
			setImageName: {deleteImage, panic},
			deleteImage:  {completed, panic},
		},
		[]telegram.State{completed, panic},
	)
}

func (u *deleteImageBox) setImageNameAction(fsmContex telegram.FSMContext, event telegram.UserEvent) telegram.StateResult {
	u.name = event.RawMessage.Text
	return usecases.OkState(deleteImage)
}

func (u *deleteImageBox) deleteImage(fsmContex telegram.FSMContext, _ <-chan telegram.UserEvent) telegram.StateResult {
	if err := u.repo.DeleteByName(u.name); err != nil {
		return usecases.ErrorState(panic, err)
	}
	return usecases.OkState(completed)
}
