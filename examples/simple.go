package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	state_go "github.com/inamvar/state-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MyPayload struct {
	Counter int    `json:"counter"`
	Message string `json:"message"`
}

func exampleCallback(ctx context.Context, payload interface{}) (state_go.Status, interface{}, error) {
	p, ok := payload.(*MyPayload)
	if !ok {
		return state_go.StatusFailed, nil, errors.New("invalid payload type")
	}

	p.Counter++
	fmt.Printf("Processing: Counter=%d, Message=%s\n", p.Counter, p.Message)
	//
	//if p.Counter < 3 {
	//	return state_go.StatusUnknown, p, nil // retry current state
	//}
	//if p.Message == "fail" {
	//	return state_go.StatusFailed, p, nil
	//}

	return state_go.StatusSuccess, p, nil
}

func main() {
	dsn := "host=localhost user=postgres password=boofhichkas dbname=statemachine port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Migrate schema
	if err := db.AutoMigrate(&state_go.Job{}); err != nil {
		panic(err)
	}

	sm := state_go.NewStateMachine(db, 3)

	stateA := state_go.State{
		Name:        "state_a",
		ActionFunc:  exampleCallback,
		Transitions: map[state_go.Status]string{"success": "state_b"},
	}

	stateB := state_go.State{
		Name:        "state_b",
		ActionFunc:  exampleCallback,
		Transitions: map[state_go.Status]string{},
	}

	sm.RegisterState(stateA)
	sm.RegisterState(stateB)

	initialPayload := MyPayload{Counter: 0, Message: "hello"}

	payloadJSON, err := json.Marshal(initialPayload)
	if err != nil {
		panic(err)
	}

	job := state_go.Job{
		CurrentState: "state_a",
		Payload:      state_go.Payload(payloadJSON),
	}

	if err = db.Create(&job).Error; err != nil {
		panic(err)
	}
	ctx := context.Background()

	err = sm.Run(ctx, job.ID, &MyPayload{})
	if err != nil {
		fmt.Println("Run error:", err)
	}

}
