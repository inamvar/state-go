package statemachine

import (
	"errors"
	"fmt"
)

type State string
type Event string

type Context interface{}

type TransitionFunc func(ctx Context) (State, error)
type GuardFunc func(ctx Context) bool

type StateMachine struct {
	currentState State
	transitions  map[State]map[Event]Transition
	context      Context
	history      []State // Optional: for debugging and auditing
}

type Transition struct {
	Fn    TransitionFunc
	Guard GuardFunc
}

func NewStateMachine(initialState State) *StateMachine {
	return &StateMachine{
		currentState: initialState,
		transitions:  make(map[State]map[Event]Transition),
		context:      nil,
		history:      []State{},
	}
}

func (sm *StateMachine) AddTransition(fromState State, event Event, fn TransitionFunc, guard GuardFunc) {
	if sm.transitions[fromState] == nil {
		sm.transitions[fromState] = make(map[Event]Transition)
	}
	sm.transitions[fromState][event] = Transition{Fn: fn, Guard: guard}
}

func (sm *StateMachine) SetContext(ctx Context) {
	sm.context = ctx
}

func (sm *StateMachine) HandleEvent(event Event) error {
	transitions := sm.transitions[sm.currentState]
	if transitions == nil {
		return fmt.Errorf("no transition defined for state %s and event %s", sm.currentState, event)
	}
	transition := transitions[event]
	if transition.Fn == nil {
		return fmt.Errorf("no transition defined for state %s and event %s", sm.currentState, event)
	}
	if transition.Guard != nil && !transition.Guard(sm.context) {
		return errors.New("guard condition failed")
	}
	newState, err := transition.Fn(sm.context)
	if err != nil {
		return err
	}
	sm.currentState = newState
	sm.history = append(sm.history, newState) // Optional: record state history
	return nil
}

// Additional Considerations and Enhancements

func (sm *StateMachine) HandleEventAsync(event Event, done chan error) {
	go func() {
		err := sm.HandleEvent(event)
		done <- err
	}()
}
