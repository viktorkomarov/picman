package upload

import (
	"errors"

	"github.com/viktorkomarov/picman/internal/api/telegram"
	"github.com/viktorkomarov/picman/internal/domain"
)

const (
	setPayloadState         telegram.State = "set_payload"
	setImageNameState       telegram.State = "set_image_name"
	saveImageState          telegram.State = "save_image"
	incorrectImageNameState telegram.State = "incorrect_image_name"
	uploadCompletedState    telegram.State = "upload_completed"
	panicState              telegram.State = "panic_state"
)

func uploadImageFSM(useCase *UploadImageUseCase) *telegram.UserCommunicationFSM {
	return telegram.NewUserCommunicationFSM(
		setPayloadState,
		map[telegram.State]telegram.StateConversation{
			setPayloadState:         telegram.NewBaseConversation(telegram.NewQuestionProvider("Загрузите файл изображения"), useCase.setPayloadAction),
			setImageNameState:       telegram.NewBaseConversation(telegram.NewQuestionProvider("Укажите название изображения"), useCase.setImageNameAction),
			saveImageState:          telegram.NewBaseConversation(telegram.SkipQuestionProvider(), useCase.saveImageAction),
			incorrectImageNameState: telegram.NewBaseConversation(telegram.NewQuestionProvider("Некорректное название файла или Такой файл уже существует, укажите новое имя"), useCase.setImageNameAction),
		},
		map[telegram.State][]telegram.State{
			setPayloadState:         {setImageNameState, panicState},
			setImageNameState:       {setImageNameState, saveImageState, panicState},
			incorrectImageNameState: {saveImageState, panicState},
			saveImageState:          {incorrectImageNameState, uploadCompletedState, panicState},
		},
		map[telegram.State]telegram.TerminalState{
			panicState:           nil,
			uploadCompletedState: nil,
		},
	)
}

func (u *UploadImageUseCase) setPayloadAction(fsmContex telegram.UseCaseContext, event telegram.UserEvent) telegram.StateResult {
	err := u.imageBuilder.SetPayload(nil)
	if err != nil {
		return toPanicState(err)
	}
	return toNextState(setImageNameState)
}

func (u *UploadImageUseCase) setImageNameAction(fsmContex telegram.UseCaseContext, event telegram.UserEvent) telegram.StateResult {
	err := u.imageBuilder.SetName("")
	if err != nil {
		return toNextState(incorrectImageNameState)
	}
	return toNextState(saveImageState)
}

func (u *UploadImageUseCase) saveImageAction(fsmContex telegram.UseCaseContext, event telegram.UserEvent) telegram.StateResult {
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
