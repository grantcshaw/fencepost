package keystore

import "time"

// ExpiryPolicy defines when a key is considered expired and must be replaced.
type ExpiryPolicy struct {
	MaxAge time.Duration
}

// DefaultExpiryPolicy returns a policy that expires keys after 90 days.
func DefaultExpiryPolicy() ExpiryPolicy {
	return ExpiryPolicy{
		MaxAge: 90 * 24 * time.Hour,
	}
}

// IsExpired reports whether the key for the given entry has exceeded the
// policy's maximum age. A zero CreatedAt is treated as never expired.
func (p ExpiryPolicy) IsExpired(entry Entry) bool {
	if entry.CreatedAt.IsZero() {
		return false
	}
	return time.Since(entry.CreatedAt) > p.MaxAge
}

// TimeUntilExpiry returns the duration remaining before the entry expires
// under this policy. A negative value means the entry is already expired.
// If CreatedAt is zero, the maximum duration is returned.
func (p ExpiryPolicy) TimeUntilExpiry(entry Entry) time.Duration {
	if entry.CreatedAt.IsZero() {
		return time.Duration(1<<63 - 1)
	}
	return p.MaxAge - time.Since(entry.CreatedAt)
}

// ExpiredKeys returns the service names whose keys have exceeded the policy
// maximum age.
func ExpiredKeys(store *Store, policy ExpiryPolicy) []string {
	store.mu.RLock()
	defer store.mu.RUnlock()

	var expired []string
	for name, entry := range store.data {
		if policy.IsExpired(entry) {
			expired = append(expired, name)
		}
	}
	return expired
}
