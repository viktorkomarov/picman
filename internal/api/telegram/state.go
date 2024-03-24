package telegram

type State string

type StateResult struct {
	NextState State
	Error     error
}

func NewStateResult(next State, err error) StateResult {
	return StateResult{
		NextState: next,
		Error:     err,
	}
}
