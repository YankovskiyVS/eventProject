package transactionaloutbox

import "time"

type recordUnlocker struct {
	store           Store
	maxLockDuration time.Duration
}

func NewRecordUnlocker(store Store, maxLockTimeMinutes time.Duration) recordUnlocker {
	return recordUnlocker{
		maxLockDuration: maxLockTimeMinutes * time.Minute, // Convert minutes to duration
		store:           store,
	}
}

func (d recordUnlocker) UnlockExpiredMessages() error {
	// Calculate expiration time relative to current time
	expiryTime := time.Now().UTC().Add(-d.maxLockDuration)
	return d.store.ClearLocksWithDurationBeforeDate(expiryTime)
}
