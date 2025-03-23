package transactionaloutbox

import "time"

type recordUnlocker struct {
	store                   Store
	time                    time.Time
	MaxLockTimeDurationMins time.Duration
}

func newRecordUnlocker(store Store, maxLockTimeDurationMins time.Duration) recordUnlocker {
	return recordUnlocker{
		MaxLockTimeDurationMins: maxLockTimeDurationMins,
		store:                   store,
		time:                    time.Time{},
	}
}

func (d recordUnlocker) UnlockExpiredMessages() error {
	expiryTime := d.time.Time.UTC().Add(-d.MaxLockTimeDurationMins)
	clearErr := d.store.ClearLocksWithDurationBeforeDate(expiryTime)
	if clearErr != nil {
		return clearErr
	}
	return nil
}
