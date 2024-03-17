package upload

import (
	"errors"

	"github.com/viktorkomarov/picman/internal/api/telegram/usecases"
	"github.com/viktorkomarov/picman/internal/domain"
)

const (
	setPayloadState         usecases.State = "set_payload"
	setImageNameState       usecases.State = "set_image_name"
	saveImageState          usecases.State = "save_image"
	incorrectImageNameState usecases.State = "incorrect_image_name"
	uploadCompletedState    usecases.State = "upload_completed"
	panicState              usecases.State = "panic_state"
)

func uploadImageFSM(useCase *UploadImageUseCase) *usecases.UserCommunicationFSM {
	return usecases.NewUserCommunicationFSM(
		setPayloadState,
		map[usecases.State]usecases.StateConversation{
			setPayloadState:         usecases.NewBaseConversation(usecases.NewQuestionProvider("Загрузите файл изображения"), useCase.setPayloadAction),
			setImageNameState:       usecases.NewBaseConversation(usecases.NewQuestionProvider("Укажите название изображения"), useCase.setImageNameAction),
			saveImageState:          usecases.NewBaseConversation(usecases.SkipQuestionProvider(), useCase.saveImageAction),
			incorrectImageNameState: usecases.NewBaseConversation(usecases.NewQuestionProvider("Некорректное название файла или Такой файл уже существует, укажите новое имя"), useCase.setImageNameAction),
		},
		map[usecases.State][]usecases.State{
			setPayloadState:         {setImageNameState, panicState},
			setImageNameState:       {setImageNameState, saveImageState, panicState},
			incorrectImageNameState: {saveImageState, panicState},
			saveImageState:          {incorrectImageNameState, uploadCompletedState, panicState},
		},
		map[usecases.State]usecases.TerminalState{
			panicState:           nil,
			uploadCompletedState: nil,
		},
	)
}

func (u *UploadImageUseCase) setPayloadAction(fsmContex usecases.UseCaseContext, event usecases.UserEvent) usecases.StateResult {
	err := u.imageBuilder.SetPayload(nil)
	if err != nil {
		return toPanicState(err)
	}
	return toNextState(setImageNameState)
}

func (u *UploadImageUseCase) setImageNameAction(fsmContex usecases.UseCaseContext, event usecases.UserEvent) usecases.StateResult {
	err := u.imageBuilder.SetName("")
	if err != nil {
		return toNextState(incorrectImageNameState)
	}
	return toNextState(saveImageState)
}

func (u *UploadImageUseCase) saveImageAction(fsmContex usecases.UseCaseContext, event usecases.UserEvent) usecases.StateResult {
	err := u.saver.SaveImage(u.imageBuilder.Image())
	if err == nil {
		return toNextState(uploadCompletedState)
	}

	if errors.Is(err, domain.ErrImageAlreadyExists) {
		return toNextState(incorrectImageNameState)
	}
	return toPanicState(err)
}

func toPanicState(err error) usecases.StateResult {
	return usecases.NewStateResult(panicState, err.Error())
}

func toNextState(state usecases.State) usecases.StateResult {
	return usecases.NewStateResult(state, "")
}
