package telegram

import "fmt"

type FSMContext struct {
	rawData map[string]interface{}
	stages  []StateResult
}

func NewFSMContext(chatID int64) FSMContext {
	return FSMContext{
		rawData: map[string]interface{}{
			"chatID": interface{}(chatID),
		},
		stages: make([]StateResult, 0),
	}
}

func GetFromUseCaseContext[A any](useCaseContext FSMContext, key string) (A, error) {
	var zeroVal A

	data, ok := useCaseContext.rawData[key]
	if !ok {
		return zeroVal, fmt.Errorf("no data correlated with %s", key)
	}

	typedData, ok := data.(A)
	if !ok {
		return zeroVal, fmt.Errorf("error during type cast")
	}
	return typedData, nil
}

func (u FSMContext) LastState() StateResult {
	return u.stages[len(u.stages)-1] // len check ?
}

func (u FSMContext) WithPassedState(state StateResult) FSMContext {
	return FSMContext{
		stages:  append(u.stages, state),
		rawData: u.rawData,
	}
}
