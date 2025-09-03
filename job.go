package state_go

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
)

func (sm *StateMachine) Run(ctx context.Context, jobID uint, payloadType interface{}) error {

	tx := sm.db.Begin()
	var job Job
	// Lock the row for update (blocking concurrent transactions)
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&job, jobID).Error; err != nil {
		return err
	}
	defer func() {
		_ = tx.Commit()
	}()

	return sm.processState(ctx, tx, &job, payloadType)

}

func (sm *StateMachine) processState(ctx context.Context, tx *gorm.DB, job *Job, payloadType interface{}) error {
	state, ok := sm.states[job.CurrentState]
	if !ok {
		return errors.New("state not found for current state")
	}

	var payloadVal interface{}
	if payloadType != nil {
		payloadVal = cloneType(payloadType) // Get the correct type
		if err := json.Unmarshal(job.Payload, payloadVal); err != nil {
			return err
		}
	} else {
		payloadVal = nil
	}
	status, newPayload, err := state.ActionFunc(ctx, payloadVal)
	if err != nil {

		return err
	}

	// If newPayload is returned, marshal it and update the job
	if newPayload != nil {
		updatedPayloadJSON, err := json.Marshal(newPayload)
		if err != nil {
			return err
		}
		job.Payload = updatedPayloadJSON
	}

	nextState, hasTransition := state.Transitions[status]

	if hasTransition {

		job.CurrentState = nextState
		fmt.Printf("Job %d moved to state %s with status %s\n", job.ID, job.CurrentState, status)
		job.RetryCount = 0
		if err = tx.Save(&job).Error; err != nil {
			return err
		}

		pendingState, hasState := sm.states[nextState]
		if !hasState || pendingState.StateType == StateAsync {
			return nil
		}

		// If newPayload is of a different type, pass it to processState correctly
		if newPayload != nil {
			newPayloadType := reflect.TypeOf(newPayload)

			// Check if the newPayloadType is a pointer and matches the type we expect (*MyPayload)
			if newPayloadType.Kind() != reflect.Ptr {
				return fmt.Errorf("new payload is not a pointer: %v", newPayloadType)
			}

			// Recursively call processState with the correct type
			return sm.processState(ctx, tx, job, reflect.New(newPayloadType.Elem()).Interface())
		} else {
			// If no newPayload, continue with the current one
			return sm.processState(ctx, tx, job, nil)
		}
	} else {

		fmt.Printf("Job %d has reached the end state. Flow complete.\n", job.ID)
		// Optionally mark job as completed in DB or remove
		//job.CurrentState = nextState
		job.RetryCount = 0
		err = tx.Save(&job).Error
		return err
	}

}

//cloneType creates a new pointer to a value of the same type as 'val'

func cloneType(val interface{}) interface{} {
	return reflect.New(reflect.TypeOf(val).Elem()).Interface()
}

//func cloneType(val interface{}) interface{} {
//	// Create a new pointer to a value of the same type as 'val'
//	return reflect.New(reflect.TypeOf(val)).Interface()
//}
