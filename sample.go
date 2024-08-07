package main

import (
	"fmt"
	"github.com/inamvar/state-go/statemachine"
)

type MyContext struct {
	Data int
}

const (
	idle    statemachine.State = "idle"
	running statemachine.State = "running"
	stopped statemachine.State = "stopped"
)
const (
	startEvent statemachine.Event = "startEvent"
	stopEvent  statemachine.Event = "stopEvent"
)

func main() {
	sm := statemachine.NewStateMachine(idle)
	ctx := MyContext{Data: 42}

	sm.SetContext(ctx)

	sm.AddTransition(idle, startEvent, func(ctx statemachine.Context) (statemachine.State, error) {
		myCtx := ctx.(MyContext)
		fmt.Println("Starting with data:", myCtx.Data)
		return running, nil
	}, nil)

	sm.AddTransition(running, stopEvent, func(ctx statemachine.Context) (statemachine.State, error) {
		fmt.Println("Stopping")
		return idle, nil
	}, nil)

	err := sm.HandleEvent(startEvent)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = sm.HandleEvent(stopEvent)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
