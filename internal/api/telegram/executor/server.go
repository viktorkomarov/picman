package executor

import (
	"fmt"
	"sync"
	"time"

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
		dog := NewExecutorWatchdog(u.fsmProvder)
		u.watchdogs[userID] = dog
		go u.waitToDeleteSession(userID, dog.Loop(time.Minute*2))
	}
	u.watchdogs[userID].RecieveMessage(msg)
}

func (u *UsersHub) waitToDeleteSession(id int64, barrier <-chan struct{}) {
	go func() {
		<-barrier

		u.mu.Lock()
		defer u.mu.Unlock()

		fmt.Printf("delete %d user\n", id)
		delete(u.watchdogs, id)
	}()
}
