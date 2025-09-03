package state_go

import (
	"gorm.io/datatypes"
	"time"
)

type Status string

const (
	StatusPending Status = "pending"
	StatusSuccess Status = "success"
	StatusFailed  Status = "failed"
	StatusUnknown Status = "unknown"
)

type StateType string

const (
	StateSync  StateType = "sync"
	StateAsync StateType = "async"
)

type Payload datatypes.JSON

type Job struct {
	ID           uint `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CurrentState string  `gorm:"type:varchar(50);index"`
	Payload      Payload `gorm:"column:payload;type:jsonb"` // JSONB column in Postgres
	RetryCount   int
}

type State struct {
	Name        string       `json:"name"`
	ActionFunc  CallbackFunc `json:"-"`
	Transitions map[Status]string
	StateType   StateType
}
