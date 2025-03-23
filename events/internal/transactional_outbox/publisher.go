package transactionaloutbox

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Publisher encapsulates the publishing functionality of the outbox pattern
type Publisher struct {
	store Store
	time  time.Time
	uuid  uuid.UUID
}

// NewPublisher is the Publisher constructor
func NewPublisher(store Store) Publisher {
	return Publisher{
		store: store,
		time:  time.Time{},
		uuid:  uuid.UUID{},
	}
}

// Send stores the provided Message within the provided sql.Tx
func (p Publisher) Send(msg Message, tx *sql.Tx) error {
	newID := p.uuid.ID()
	record := Record{
		ID:          newID,
		Message:     msg,
		State:       PendingDelivery,
		CreatedOn:   p.time.UTC(),
		LockID:      nil,
		LockedOn:    nil,
		ProcessedOn: nil,
	}

	return p.store.AddRecordTx(record, tx)
}
