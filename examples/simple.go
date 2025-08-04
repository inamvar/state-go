package main

import (
	"encoding/json"
	"errors"
	"fmt"
	state_go "github.com/inamvar/state-go"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MyPayload struct {
	Counter int    `json:"counter"`
	Message string `json:"message"`
}

func exampleCallback(payload interface{}) (state_go.Status, interface{}, error) {
	p, ok := payload.(*MyPayload)
	if !ok {
		return state_go.StatusFailed, nil, errors.New("invalid payload type")
	}

	p.Counter++
	fmt.Printf("Processing: Counter=%d, Message=%s\n", p.Counter, p.Message)

	if p.Counter < 3 {
		return state_go.StatusUnknown, p, nil // retry current state
	}
	if p.Message == "fail" {
		return state_go.StatusFailed, p, nil
	}

	return state_go.StatusSuccess, p, nil
}

func main() {
	dsn := "host=localhost user=postgres password=yourpass dbname=statemachine port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Migrate schema
	if err := db.AutoMigrate(&state_go.Job{}); err != nil {
		panic(err)
	}

	sm := state_go.NewStateMachine(db, 3)

	sm.RegisterState(state_go.StateA, exampleCallback, map[state_go.Status]state_go.StateName{
		state_go.StatusSuccess: state_go.StateB,
		state_go.StatusFailed:  state_go.StateA,
		state_go.StatusUnknown: state_go.StateA,
	})

	sm.RegisterState(state_go.StateB, exampleCallback, map[state_go.Status]state_go.StateName{
		state_go.StatusSuccess: state_go.StateEnd,
		state_go.StatusFailed:  state_go.StateB,
		state_go.StatusUnknown: state_go.StateEnd,
	})
	initialPayload := MyPayload{Counter: 0, Message: "hello"}

	payloadJSON, err := json.Marshal(initialPayload)
	if err != nil {
		panic(err)
	}

	job := state_go.Job{
		CurrentState: state_go.StateA,
		Payload:      datatypes.JSON(payloadJSON),
	}

	if err := db.Create(&job).Error; err != nil {
		panic(err)
	}

	for i := 0; i < 5; i++ {
		err := sm.Run(job.ID, &MyPayload{})
		if err != nil {
			fmt.Println("Run error:", err)
			break
		}
	}
}
