
## statemachine

**A flexible and customizable state machine implementation in Go.**

This package provides a generic state machine framework that can be adapted to various use cases. It supports features such as:

* **Generic state and event definitions**
* **Customizable context for state transitions**
* **Guard conditions for controlling transitions**
* **Optional state history tracking**
* **Asynchronous event handling**

### Usage

```go
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

```

### API

* **NewStateMachine(initialState State)**: Creates a new state machine with the specified initial state.
* **AddTransition(fromState State, event Event, fn TransitionFunc, guard GuardFunc)**: Adds a transition from the `fromState` to a new state based on the `event`, with optional `guard` condition.
* **SetContext(ctx Context)**: Sets the context for the state machine.
* **HandleEvent(event Event)**: Handles the specified event and transitions the state machine accordingly.
* **HandleEventAsync(event Event, done chan error)**: Handles the event asynchronously.

### Additional Features

* **Hierarchical states**
* **Parallel transitions**
* **Timeout handling**
* **Entry/exit actions**
* **State history management**

### Contributions

Contributions are welcome! Please open an issue or pull request if you have any suggestions or improvements.

### License

This project is licensed under the MIT License.


