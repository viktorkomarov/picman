package upload

import (
	"context"
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
	"github.com/viktorkomarov/picman/internal/api/telegram"
	"github.com/viktorkomarov/picman/internal/api/telegram/usecases"
	"github.com/viktorkomarov/picman/internal/domain"
	"github.com/viktorkomarov/picman/internal/fetcher"
)

const (
	setPayloadState         telegram.State = "set_payload"
	setImageNameState       telegram.State = "set_image_name"
	saveImageState          telegram.State = "save_image"
	incorrectImageNameState telegram.State = "incorrect_image_name"
	uploadCompletedState    telegram.State = "upload_completed"
	panicState              telegram.State = "panic_state"
)

type uploadImageBox struct {
	saver        domain.ImageRepository
	imageBuilder *domain.ImageBuilder
	fetcher      *fetcher.TelegramImageFetcher
	rawData      map[string]interface{}
}

func NewFSM(saver domain.ImageRepository, bot *tgbotapi.BotAPI, fetcher *fetcher.TelegramImageFetcher) *telegram.FSM {
	box := &uploadImageBox{
		saver:        saver,
		imageBuilder: domain.NewImageBuilder(),
		fetcher:      fetcher,
		rawData:      make(map[string]interface{}),
	}

	return telegram.NewFSM(
		setPayloadState,
		map[telegram.State]telegram.StateAction{
			setPayloadState: usecases.NewStateAction(
				usecases.SendMessageNotifyFunc(bot, "Загрузите файл изображения"),
				usecases.ActionWithEvent(panicState, box.setPayloadAction),
			),
			setImageNameState: usecases.NewStateAction(
				usecases.SendMessageNotifyFunc(bot, "Укажите название изображения [произвольное имя файл].[png|jpg|gif]"),
				usecases.ActionWithEvent(panicState, box.setImageNameAction),
			),
			saveImageState: usecases.NewStateAction(
				usecases.EmptyNotifyFunc(), box.saveImageAction,
			),
			incorrectImageNameState: usecases.NewStateAction(
				usecases.SendMessageNotifyFunc(bot, "Некорректное название файла или Такой файл уже существует, укажите новое имя"),
				usecases.ActionWithEvent(panicState, box.setImageNameAction),
			),
			uploadCompletedState: usecases.NewStateAction(
				usecases.SendMessageNotifyFunc(bot, "Файл успешно загружен!"),
				usecases.EmptyAction(),
			),
			panicState: usecases.NewStateAction(usecases.ErrorUserNotify(bot), usecases.EmptyAction()),
		},
		map[telegram.State][]telegram.State{
			setPayloadState:         {setImageNameState, panicState},
			setImageNameState:       {incorrectImageNameState, saveImageState, panicState},
			incorrectImageNameState: {incorrectImageNameState, saveImageState, panicState},
			saveImageState:          {incorrectImageNameState, uploadCompletedState, panicState},
		},
		[]telegram.State{panicState, uploadCompletedState},
	)
}

func (u *uploadImageBox) setPayloadAction(fsmContex telegram.FSMContext, event telegram.UserEvent) telegram.StateResult {
	if len(event.RawMessage.Photo) == 0 {
		return usecases.ErrorState(panicState, fmt.Errorf("photo should be presented"))
	}
	maxQualityPhoto := lo.MaxBy(event.RawMessage.Photo, func(lhs, rhs tgbotapi.PhotoSize) bool {
		return lhs.FileSize > rhs.FileSize
	})
	payload, err := u.fetcher.Fetch(context.Background(), maxQualityPhoto.FileID)
	if err != nil {
		return usecases.ErrorState(panicState, err)
	}
	if err = u.imageBuilder.SetPayload(payload); err != nil {
		return usecases.ErrorState(panicState, err)
	}
	return usecases.OkState(setImageNameState)
}

func (u *uploadImageBox) setImageNameAction(fsmContex telegram.FSMContext, event telegram.UserEvent) telegram.StateResult {
	err := u.imageBuilder.SetName(event.RawMessage.Text)
	if err != nil {
		return usecases.OkState(incorrectImageNameState)
	}
	return usecases.OkState(saveImageState)
}

func (u *uploadImageBox) saveImageAction(fsmContex telegram.FSMContext, _ <-chan telegram.UserEvent) telegram.StateResult {
	err := u.saver.SaveImage(u.imageBuilder.Image())
	if err == nil {
		return usecases.OkState(uploadCompletedState)
	}

	if errors.Is(err, domain.ErrImageAlreadyExists) {
		return usecases.OkState(incorrectImageNameState)
	}
	return usecases.ErrorState(panicState, err)
}
