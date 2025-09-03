package state_go

import (
	"context"
	"gorm.io/gorm"
	"sync"
)

type CallbackFunc func(ctx context.Context, payload interface{}) (Status, interface{}, error)

type StateMachine struct {
	db         *gorm.DB
	states     map[string]State
	maxRetries int
	mu         sync.RWMutex
}

func NewStateMachine(db *gorm.DB, maxRetries int) *StateMachine {
	return &StateMachine{
		db:         db,
		states:     make(map[string]State),
		maxRetries: maxRetries,
		mu:         sync.RWMutex{},
	}
}

func (sm *StateMachine) RegisterState(state State) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.states[state.Name] = state

}
