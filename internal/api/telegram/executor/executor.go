package executor

import (
	"errors"
	"fmt"

	"github.com/viktorkomarov/picman/internal/api/telegram"
)

type UseCaseExecution struct {
	context telegram.FSMContext
	fsm     *telegram.FSM
	msg     <-chan telegram.UserEvent
}

func NewUserExecution(fsm *telegram.FSM, msg <-chan telegram.UserEvent, chatID int64) *UseCaseExecution {
	execution := &UseCaseExecution{
		fsm:     fsm,
		msg:     msg,
		context: telegram.NewFSMContext(chatID),
	}

	return execution
}

func (u *UseCaseExecution) Run() <-chan error {
	errCh := make(chan error)

	go func() {
		defer close(errCh)

		for {
			if err := u.fsm.NotifyUser(u.context); err != nil {
				errCh <- fmt.Errorf("notifyUser: %w", err)
				return
			}

			nextState := u.fsm.ApplyUserEvent(u.context, u.msg)
			u.context = u.context.WithPassedState(nextState)

			err := u.fsm.Transit(nextState)
			switch {
			case errors.Is(err, telegram.ErrEndOfFSM):
				return
			case err != nil:
				errCh <- fmt.Errorf("transit: %w", err)
				return
			}
		}
	}()

	return errCh
}
