package telegram

import (
	"github.com/samber/lo"
)

type State string

type StateConversation interface {
	ApplyEvent(UseCaseContext, UserEvent) StateResult
	Question() QuestionProvider
}

type TerminalState interface {
	Output() Output
}

type UserCommunicationFSM struct {
	current        State
	actions        map[State]StateConversation
	transitions    map[State]map[State]bool
	terminalStates map[State]TerminalState
}

func statesToSet(states []State) map[State]bool {
	return lo.SliceToMap(states, func(state State) (State, bool) {
		return state, true
	})
}

func NewUserCommunicationFSM(
	init State,
	actions map[State]StateConversation,
	transitions map[State][]State,
	terminalStates map[State]TerminalState,
) *UserCommunicationFSM {
	transitionSet := make(map[State]map[State]bool)
	for state := range transitions {
		transitionSet[state] = statesToSet(transitions[state])
	}

	return &UserCommunicationFSM{
		current:        init,
		actions:        actions,
		transitions:    transitionSet,
		terminalStates: terminalStates,
	}
}

func (u *UserCommunicationFSM) CurrentQuestion() QuestionProvider {
	return u.actions[u.current].Question()
}

func (u *UserCommunicationFSM) Transit(fsmContext UseCaseContext, event UserEvent) UseCaseContext {
	action := u.actions[u.current]

	stateResult := action.ApplyEvent(fsmContext, event)

	currentTransitions := u.transitions[stateResult.NextState]
	if !currentTransitions[stateResult.NextState] {
		panic("miconfigured fsm")
	}

	u.current = stateResult.NextState
	// add step
	return fsmContext
}

func (u *UserCommunicationFSM) IsTerminalState() bool {
	_, ok := u.terminalStates[u.current]
	return ok
}

func (u *UserCommunicationFSM) Terminal() Output {
	return u.terminalStates[u.current].Output()
}
