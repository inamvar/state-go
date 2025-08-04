package state_go

import (
	"gorm.io/gorm"
)

type CallbackFunc func(payload interface{}) (Status, interface{}, error)

type StateMachine struct {
	db          *gorm.DB
	callbacks   map[StateName]CallbackFunc
	transitions map[StateName]map[Status]StateName
	maxRetries  int
}

func NewStateMachine(db *gorm.DB, maxRetries int) *StateMachine {
	return &StateMachine{
		db:          db,
		callbacks:   make(map[StateName]CallbackFunc),
		transitions: make(map[StateName]map[Status]StateName),
		maxRetries:  maxRetries,
	}
}

func (sm *StateMachine) RegisterState(state StateName, callback CallbackFunc, transitions map[Status]StateName) {
	sm.callbacks[state] = callback
	sm.transitions[state] = transitions
}
