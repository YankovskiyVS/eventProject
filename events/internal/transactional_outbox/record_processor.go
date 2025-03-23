package transactionaloutbox

import (
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

			if updateErr := d.store.UpdateRecordByID(rec); updateErr != nil {
				finalErr := fmt.Errorf("multiple errors occurred: %w while handling %v", updateErr, finalErr)
				return finalErr
			}

			finalErr = fmt.Errorf("failed to send message %s: %w (record update attempted)", rec.ID, err)
			continue
		}

		// Success case
		rec.State = Delivered
		rec.LockedOn = nil
		rec.LockID = nil
		rec.ProcessedOn = &now
		if err := d.store.UpdateRecordByID(rec); err != nil {
			finalErr = fmt.Errorf("failed to update delivered record %s: %w", rec.ID, err)
		}
	}
	return finalErr
}

// lockUnprocessedEntities updates the messages with the current machine's lockID
func (d defaultRecordProcessor) lockUnprocessedEntities() error {
	lockTime := time.Now().UTC()
	return d.store.UpdateRecordLockByState(d.machineID, lockTime, PendingDelivery)
}
