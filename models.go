package state_go

import (
	"gorm.io/datatypes"
	"time"
)

type Status string

const (
	StatusSuccess Status = "success"
	StatusFailed  Status = "failed"
	StatusUnknown Status = "unknown"
)

type StateName string

const (
	StateA   StateName = "StateA"
	StateB   StateName = "StateB"
	StateC   StateName = "StateC"
	StateEnd StateName = "END"
	// add more states as needed
)

type Job struct {
	ID           uint `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CurrentState StateName      `gorm:"type:varchar(50);index"`
	Payload      datatypes.JSON `gorm:"type:jsonb"` // JSONB column in Postgres
	RetryCount   int
}
