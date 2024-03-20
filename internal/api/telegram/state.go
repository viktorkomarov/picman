package telegram

type State string

type StateResult struct {
	NextState    State
	FailedReason string
}

func NewStateResult(next State, failedReason string) StateResult {
	return StateResult{
		NextState:    next,
		FailedReason: failedReason,
	}
}
