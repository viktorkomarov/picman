package telegram

import (
	"errors"
	"fmt"

	"github.com/samber/lo"
)

type StateAction interface {
	NotifyUser(FSMContext) error
	ApplyUserEvent(FSMContext, UserEvent) StateResult
}

var (
	ErrEndOfFSM         = errors.New("fsm is executed")
	ErrMisconfiguredFSM = errors.New("misconfigured fsm")
)

type FSMType int

const (
	FSMTypeUpload FSMType = iota
)

type FSM struct {
	current        State
	actions        map[State]StateAction
	transitions    map[State]map[State]bool
	terminalStates map[State]bool
}

func statesToSet(states []State) map[State]bool {
	return lo.SliceToMap(states, func(state State) (State, bool) {
		return state, true
	})
}

func NewFSM(
	init State,
	actions map[State]StateAction,
	transitions map[State][]State,
	terminalStates []State,
) *FSM {
	transitionSet := make(map[State]map[State]bool)
	for state := range transitions {
		transitionSet[state] = statesToSet(transitions[state])
	}

	return &FSM{
		current:        init,
		actions:        actions,
		transitions:    transitionSet,
		terminalStates: statesToSet(terminalStates),
	}
}

func (f *FSM) NotifyUser(ctx FSMContext) error {
	return f.actions[f.current].NotifyUser(ctx)
}

func (f *FSM) ApplyUserEvent(ctx FSMContext, event UserEvent) StateResult {
	return f.actions[f.current].ApplyUserEvent(ctx, event)
}

func (f *FSM) Transit(state StateResult) error {
	if f.terminalStates[f.current] {
		return ErrEndOfFSM
	}

	_, ok := f.transitions[f.current][state.NextState]
	if !ok {
		return fmt.Errorf("no transition from %s to %s: %w", f.current, state.NextState, ErrMisconfiguredFSM)
	}
	f.current = state.NextState

	return nil
}
