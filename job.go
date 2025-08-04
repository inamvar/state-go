package state_go

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
)

func (sm *StateMachine) Run(jobID uint, payloadType interface{}) error {
	return sm.db.Transaction(func(tx *gorm.DB) error {
		var job Job
		// Lock the row for update (blocking concurrent transactions)
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&job, jobID).Error; err != nil {
			return err
		}

		callback, ok := sm.callbacks[job.CurrentState]
		if !ok {
			return errors.New("callback not found for current state")
		}

		payloadVal := cloneType(payloadType)
		if err := json.Unmarshal(job.Payload, payloadVal); err != nil {
			return err
		}

		status, newPayload, err := callback(payloadVal)
		if err != nil {
			return err
		}

		if newPayload != nil {
			updatedPayloadJSON, err := json.Marshal(newPayload)
			if err != nil {
				return err
			}
			job.Payload = datatypes.JSON(updatedPayloadJSON)
		}

		nextState, hasTransition := sm.transitions[job.CurrentState][status]

		if hasTransition {
			if nextState == StateEnd {
				fmt.Printf("Job %d has reached the end state. Flow complete.\n", job.ID)
				// Optionally mark job as completed in DB or remove
				job.CurrentState = nextState
				job.RetryCount = 0
				if err := tx.Save(&job).Error; err != nil {
					return err
				}
				return nil // end the flow here
			}
			job.CurrentState = nextState
			job.RetryCount = 0
		} else {
			if job.RetryCount < sm.maxRetries {
				job.RetryCount++
			} else {
				return errors.New("max retry attempts reached for state " + string(job.CurrentState))
			}
		}

		if err := tx.Save(&job).Error; err != nil {
			return err
		}

		fmt.Printf("Job %d moved to state %s with status %s\n", job.ID, job.CurrentState, status)
		return nil
	})
}

// cloneType creates a new pointer to a value of the same type as 'val'
func cloneType(val interface{}) interface{} {
	return reflect.New(reflect.TypeOf(val).Elem()).Interface()
}
