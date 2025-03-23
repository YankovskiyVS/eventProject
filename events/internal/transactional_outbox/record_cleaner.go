package transactionaloutbox

import "time"

type recordCleaner struct {
	store             Store
	time              time.Time
	MaxRecordLifetime time.Duration
}

func newRecordCleaner(store Store, maxRecordLifetime time.Duration) recordCleaner {
	return recordCleaner{
		MaxRecordLifetime: maxRecordLifetime,
		store:             store,
		time:              time.Time{},
	}
}

func (d recordCleaner) RemoveExpiredMessages() error {
	expiryTime := d.time.Time.UTC().Add(-d.MaxRecordLifetime)
	err := d.store.RemoveRecordsBeforeDatetime(expiryTime)
	if err != nil {
		return err
	}
	return nil
}
