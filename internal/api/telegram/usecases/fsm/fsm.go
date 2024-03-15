package fsm

import (
	"fmt"

	"github.com/samber/lo"
	"github.com/viktorkomarov/picman/internal/api/telegram"
)

type State string

type StateConversation interface {
	ApplyEvent(event telegram.UserEvent) (State, error)
	Question() telegram.Question
}

type UserCommunicationFSM struct {
	current        State
	actions        map[State]StateConversation
	transitions    map[State]map[State]bool
	terminalStates map[State]bool
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
	terminalStates []State,
) *UserCommunicationFSM {
	transitionSet := make(map[State]map[State]bool)
	for state := range transitions {
		transitionSet[state] = statesToSet(transitions[state])
	}

	return &UserCommunicationFSM{
		current:        init,
		actions:        actions,
		transitions:    transitionSet,
		terminalStates: statesToSet(terminalStates),
	}
}

func (u *UserCommunicationFSM) CurrentQuestion() telegram.Question {
	return u.actions[u.current].Question()
}

func (u *UserCommunicationFSM) Transit(event telegram.UserEvent) error {
	action, ok := u.actions[u.current]
	if !ok {
		return fmt.Errorf("miconfigured fsm")
	}

	nextState, err := action.ApplyEvent(event)
	if err != nil {
		return err
	}

	currentTransitions := u.transitions[nextState]
	if !currentTransitions[nextState] {
		return fmt.Errorf("miconfigured fsm")
	}

	u.current = nextState
	return nil
}

func (u *UserCommunicationFSM) IsTerminalState() bool {
	return u.terminalStates[u.current]
}
