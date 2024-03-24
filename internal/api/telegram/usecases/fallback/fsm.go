package fallback

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/viktorkomarov/picman/internal/api/telegram"
	"github.com/viktorkomarov/picman/internal/api/telegram/usecases"
)

const (
	sendList  telegram.State = "send_wellcome_message"
	completed telegram.State = "completed"
	panic     telegram.State = "panic"
)

const wellcomeMessage = `
Я умею выполнять следующие команды:
/upload_image - загрузить фотографию
/get_by_name - получить фотогграфию по имени
/list_images_name - получить имена сохраненных фотографий
/delete_by_name - удалить файл по имени
`

func NewFSM(bot *tgbotapi.BotAPI) *telegram.FSM {
	return telegram.NewFSM(
		sendList,
		map[telegram.State]telegram.StateAction{
			sendList: usecases.NewStateAction(
				usecases.SendMessageNotifyFunc(bot, wellcomeMessage),
				usecases.EmptyAction(),
			),
		},
		map[telegram.State][]telegram.State{},
		[]telegram.State{sendList},
	)
}
