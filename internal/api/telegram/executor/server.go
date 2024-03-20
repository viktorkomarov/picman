package executor

import (
	"sync"

	"github.com/viktorkomarov/picman/internal/api/telegram"
	"github.com/viktorkomarov/picman/internal/api/telegram/usecases/provider"
)

type UsersHub struct {
	mu         sync.Mutex
	watchdogs  map[int64]*ExecutorWatchdog
	fsmProvder *provider.FSMBuilder
}

func NewUserHub(fsmProvider *provider.FSMBuilder) *UsersHub {
	return &UsersHub{
		watchdogs:  make(map[int64]*ExecutorWatchdog),
		fsmProvder: fsmProvider,
	}
}

func (u *UsersHub) OnRecieveUserMessage(msg telegram.UserEvent) {
	u.mu.Lock()
	defer u.mu.Unlock()

	userID := msg.RawMessage.From.ID
	if _, ok := u.watchdogs[userID]; !ok {
		u.watchdogs[userID] = NewExecutorWatchdog(u.fsmProvder)
	}
	u.watchdogs[userID].RecieveMessage(msg)
}
