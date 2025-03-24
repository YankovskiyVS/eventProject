package transactionaloutbox

import (
	"time"
)

type recordCleaner struct {
	store             Store
	MaxRecordLifetime time.Duration
}

func NewRecordCleaner(store Store, maxRecordLifetime time.Duration) recordCleaner {
	return recordCleaner{
		MaxRecordLifetime: maxRecordLifetime,
		store:             store,
	}
}

func (d recordCleaner) RemoveExpiredMessages() error {
	expiryTime := time.Now().UTC().Add(-d.MaxRecordLifetime)
	return d.store.RemoveRecordsBeforeDatetime(expiryTime)
}
