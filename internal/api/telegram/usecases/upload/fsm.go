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

func NewUploadImageFSM(saver domain.ImageRepository, bot *tgbotapi.BotAPI, fetcher *fetcher.TelegramImageFetcher) *telegram.FSM {
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
				usecases.SendToUserMessage(bot, "Загрузите файл изображения"),
				box.setPayloadAction,
			),
			setImageNameState: usecases.NewStateAction(
				usecases.SendToUserMessage(bot, "Укажите название изображения [произвольное имя файл].[png|jpg|gif]"),
				box.setImageNameAction,
			),
			saveImageState: usecases.NewStateAction(
				usecases.EmptyNotifyFunc(), box.saveImageAction,
			),
			incorrectImageNameState: usecases.NewStateAction(
				usecases.SendToUserMessage(bot, "Некорректное название файла или Такой файл уже существует, укажите новое имя"),
				box.setImageNameAction,
			),
			uploadCompletedState: usecases.NewStateAction(
				usecases.SendToUserMessage(bot, "Файл успешно загружен!"),
				usecases.EmptyAction(),
			),
			panicState: usecases.NewStateAction(
				usecases.SendToUserMessage(bot, "Упс, возникли какие-то трудности - повтори попытку позже"),
				usecases.EmptyAction(),
			),
		},
		map[telegram.State][]telegram.State{
			setPayloadState:         {setImageNameState, panicState},
			setImageNameState:       {setImageNameState, saveImageState, panicState},
			incorrectImageNameState: {saveImageState, panicState},
			saveImageState:          {incorrectImageNameState, uploadCompletedState, panicState},
		},
		[]telegram.State{panicState, uploadCompletedState},
	)
}

func (u *uploadImageBox) setPayloadAction(fsmContex telegram.FSMContext, eventCh <-chan telegram.UserEvent) telegram.StateResult {
	event, ok := <-eventCh
	if !ok {
		return toPanicState(fmt.Errorf("expected to receive user event"))
	}
	if len(event.RawMessage.Photo) == 0 {
		return toPanicState(fmt.Errorf("photo should be presented"))
	}
	maxQualityPhoto := lo.MaxBy(event.RawMessage.Photo, func(lhs, rhs tgbotapi.PhotoSize) bool {
		return lhs.FileSize > rhs.FileSize
	})
	payload, err := u.fetcher.Fetch(context.Background(), maxQualityPhoto.FileID)
	if err != nil {
		return toPanicState(err)
	}
	if err = u.imageBuilder.SetPayload(payload); err != nil {
		return toPanicState(err)
	}
	return toNextState(setImageNameState)
}

func (u *uploadImageBox) setImageNameAction(fsmContex telegram.FSMContext, eventCh <-chan telegram.UserEvent) telegram.StateResult {
	event, ok := <-eventCh
	if !ok {
		return toPanicState(fmt.Errorf("expected to receive user event"))
	}

	err := u.imageBuilder.SetName(event.RawMessage.Text)
	if err != nil {
		return toNextState(incorrectImageNameState)
	}
	return toNextState(saveImageState)
}

func (u *uploadImageBox) saveImageAction(fsmContex telegram.FSMContext, _ <-chan telegram.UserEvent) telegram.StateResult {
	err := u.saver.SaveImage(u.imageBuilder.Image())
	if err == nil {
		return toNextState(uploadCompletedState)
	}

	if errors.Is(err, domain.ErrImageAlreadyExists) {
		return toNextState(incorrectImageNameState)
	}
	return toPanicState(err)
}

func toPanicState(err error) telegram.StateResult {
	return telegram.NewStateResult(panicState, err.Error())
}

func toNextState(state telegram.State) telegram.StateResult {
	return telegram.NewStateResult(state, "")
}
