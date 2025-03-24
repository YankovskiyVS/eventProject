package transactionaloutbox

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Publisher encapsulates the publishing functionality of the outbox pattern
type Publisher struct {
	store Store
}

// NewPublisher is the Publisher constructor
func NewPublisher(store Store) Publisher {
	return Publisher{
		store: store,
	}
}

// Send stores the provided Message within the provided sql.Tx
func (p Publisher) Send(msg Message, tx *sql.Tx) error {
	newID := uuid.New()
	record := Record{
		ID:          newID,
		Message:     msg,
		State:       PendingDelivery,
		CreatedOn:   time.Now().UTC(),
		LockID:      nil,
		LockedOn:    nil,
		ProcessedOn: nil,
	}

	return p.store.AddRecordTx(record, tx)
}
