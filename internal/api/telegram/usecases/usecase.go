package usecases

import "github.com/viktorkomarov/picman/internal/api/telegram"

type UserEvent struct {
	Msg telegram.Message
}

type StateResult struct {
	NextState    State
	FailedReason string
}

func NewStateResult(nextState State, reason string) StateResult {
	return StateResult{
		NextState:    nextState,
		FailedReason: reason,
	}
}

type UseCaseContext struct {
	stages []StateResult
}

func (u UseCaseContext) WithPassedState(state StateResult) UseCaseContext {
	return UseCaseContext{
		stages: append(u.stages, state),
	}
}

type Output struct{}

type UseCase interface {
	Next() bool
	Question() QuestionProvider
	ApplyUserEvent(UseCaseContext, UserEvent) StateResult
	Output() Output
}
