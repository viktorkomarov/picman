package executor

import (
	"github.com/viktorkomarov/picman/internal/api/telegram"
	"github.com/viktorkomarov/picman/internal/api/telegram/usecases"
)

type UseCaseExecution struct {
	userHandler *UsersHub
	useCase     usecases.UseCase
	msg         <-chan telegram.Message
}

func NewUserExecution(userHandler *UsersHub, useCase usecases.UseCase, msg <-chan telegram.Message) *UseCaseExecution {
	execution := &UseCaseExecution{
		userHandler: userHandler,
		useCase:     useCase,
	}
	execution.run()

	return execution
}

func (u *UseCaseExecution) run() {
	var execContext usecases.UseCaseContext

	for u.useCase.Next() {
		provider := u.useCase.Question()
		if !provider.ShouldAsk() {
			/*if err := u.userHandler.SendMessage(); err != nil {
				// toPanicState
			}*/
		}

		msg, ok := <-u.msg
		if !ok {
			return
		}

		passed := u.useCase.ApplyUserEvent(execContext, usecases.UserEvent{Msg: msg})
		execContext = execContext.WithPassedState(passed)
	}

	// make output
}
