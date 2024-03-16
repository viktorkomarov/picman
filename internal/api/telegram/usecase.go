package telegram

type UserEvent struct {
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

type Output struct{}

type UseCase interface {
	Next() bool
	Question() QuestionProvider
	ApplyUserEvent(UseCaseContext, UserEvent)
	Output() Output
}
