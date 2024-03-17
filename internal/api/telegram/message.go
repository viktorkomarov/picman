package telegram

type Message struct {
	UserID int64
}

func (m Message) IsCommand() bool {
	return false
}
