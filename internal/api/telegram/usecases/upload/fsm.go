package upload

import (
	"github.com/viktorkomarov/picman/internal/api/telegram"
	"github.com/viktorkomarov/picman/internal/api/telegram/usecases/fsm"
)

const (
	setPayloadState      fsm.State = "set_payload"
	setImageNameState    fsm.State = "set_image_name"
	saveImageState       fsm.State = "save_image"
	uploadCompletedState fsm.State = "upload_completed"
	panicState           fsm.State = "panic_state"
)

func uploadImageFSM() *fsm.UserCommunicationFSM {
	return fsm.NewUserCommunicationFSM(
		setPayloadState,
		map[fsm.State]fsm.StateConversation{
			setPayloadState:      setPayloadStateConversation{},
			setImageNameState:    setImageNameStateConversation{},
			saveImageState:       saveImageStateConversation{},
			uploadCompletedState: uploadCompletedStateConversation{},
			panicState:           panicStateConversation{},
		},
		map[fsm.State][]fsm.State{
			setPayloadState:   {setImageNameState, panicState},
			setImageNameState: {saveImageState, panicState},
			saveImageState:    {uploadCompletedState, panicState},
		},
		[]fsm.State{uploadCompletedState, panicState},
	)
}

type setPayloadStateConversation struct{}

func (s setPayloadStateConversation) Question() telegram.Question {
	return telegram.Question{}
}

func (s setPayloadStateConversation) ApplyEvent(event telegram.UserEvent) (fsm.State, error) {
	return setImageNameState, nil
}

type setImageNameStateConversation struct{}

func (s setImageNameStateConversation) Question() telegram.Question {
	return telegram.Question{}
}

func (s setImageNameStateConversation) ApplyEvent(event telegram.UserEvent) (fsm.State, error) {
	return setImageNameState, nil
}

type saveImageStateConversation struct{}

func (s saveImageStateConversation) Question() telegram.Question {
	return telegram.Question{}
}

func (s saveImageStateConversation) ApplyEvent(event telegram.UserEvent) (fsm.State, error) {
	return setImageNameState, nil
}

type uploadCompletedStateConversation struct{}

func (s uploadCompletedStateConversation) Question() telegram.Question {
	return telegram.Question{}
}

func (s uploadCompletedStateConversation) ApplyEvent(event telegram.UserEvent) (fsm.State, error) {
	return setImageNameState, nil
}

type panicStateConversation struct{}

func (s panicStateConversation) Question() telegram.Question {
	return telegram.Question{}
}

func (s panicStateConversation) ApplyEvent(event telegram.UserEvent) (fsm.State, error) {
	return setImageNameState, nil
}
