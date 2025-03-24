package transactionaloutbox

import (
	"errors"
	"fmt"
	"time"
)

// defaultRecordProcessor checks and dispatches new messages to be sent
type defaultRecordProcessor struct {
	messageBroker MessageBroker
	store         Store
	machineID     string
	retrialPolicy RetrialPolicy
}

// NewProcessor constructs a new defaultRecordProcessor
func NewProcessor(store Store, messageBroker MessageBroker, machineID string, retrialPolicy RetrialPolicy) *defaultRecordProcessor {
	return &defaultRecordProcessor{
		messageBroker: messageBroker,
		store:         store,
		machineID:     machineID,
		retrialPolicy: retrialPolicy,
	}
}

// ProcessRecords locks unprocessed messages, tries to deliver them and then unlocks them
func (d defaultRecordProcessor) ProcessRecords() error {
	if err := d.lockUnprocessedEntities(); err != nil {
		return err
	}
	defer d.store.ClearLocksByLockID(d.machineID)

	records, err := d.store.GetRecordsByLockID(d.machineID)
	if err != nil {
		return err
	}

	return d.publishMessages(records)
}

func (d defaultRecordProcessor) publishMessages(records []Record) error {
	var finalErr error
	for _, rec := range records {
		now := time.Now().UTC()
		rec.LastAttemptOn = &now
		rec.NumberOfAttempts++

		err := d.messageBroker.Send(rec.Message)
		if err != nil {
			rec.LockedOn = nil
			rec.LockID = nil
			errorMsg := err.Error()
			rec.Error = &errorMsg

			if d.retrialPolicy.MaxSendAttemptsEnabled && rec.NumberOfAttempts >= d.retrialPolicy.MaxSendAttempts {
				rec.State = MaxAttemptsReached
			}

			// Handle update error and accumulate both errors
			if updateErr := d.store.UpdateRecordByID(rec); updateErr != nil {
				finalErr = errors.Join(finalErr,
					fmt.Errorf("update failed for %s: %w", rec.ID, updateErr))
			}
			// Accumulate send error
			finalErr = errors.Join(finalErr,
				fmt.Errorf("send failed for %s: %w", rec.ID, err))
			continue
		}

		// Success case handling
		rec.State = Delivered
		rec.LockedOn = nil
		rec.LockID = nil
		rec.ProcessedOn = &now
		if err := d.store.UpdateRecordByID(rec); err != nil {
			// Accumulate delivery update error
			finalErr = errors.Join(finalErr,
				fmt.Errorf("delivery update failed for %s: %w", rec.ID, err))
		}
	}
	return finalErr
}

// lockUnprocessedEntities updates the messages with the current machine's lockID
func (d defaultRecordProcessor) lockUnprocessedEntities() error {
	lockTime := time.Now().UTC()
	return d.store.UpdateRecordLockByState(d.machineID, lockTime, PendingDelivery)
}
