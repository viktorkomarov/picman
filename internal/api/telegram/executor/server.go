package executor

import (
	"sync"

	"github.com/viktorkomarov/picman/internal/api/telegram"
)

type UsersHub struct {
	mu        sync.Mutex
	watchdogs map[int64]*ExecutorWatchdog
}

func (u *UsersHub) OnRecieveUserMessage(msg telegram.Message) {
	u.mu.Lock()
	defer u.mu.Unlock()

	if _, ok := u.watchdogs[msg.UserID]; !ok {
		u.watchdogs[msg.UserID] = NewExecutorWatchdog(nil)
	}
	u.watchdogs[msg.UserID].RecieveMessage(msg)
}
