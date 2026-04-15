package keystore

import "time"

// KeyStatus summarizes the health and state of a single service key.
type KeyStatus struct {
	Service     string
	CreatedAt   time.Time
	RotatedAt   time.Time
	IsExpired   bool
	DueRotation bool
	Tags        []string
	Note        string
}

// Status returns a KeyStatus for the named service.
// Returns an error if the service does not exist.
func (s *Store) Status(service string) (KeyStatus, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Keys[service]
	if !ok {
		return KeyStatus{}, ErrNotFound
	}

	expiry := DefaultExpiryPolicy()
	rotation := DefaultRotationPolicy()

	return KeyStatus{
		Service:     service,
		CreatedAt:   entry.CreatedAt,
		RotatedAt:   entry.RotatedAt,
		IsExpired:   expiry.IsExpired(entry),
		DueRotation: rotation.DueForRotation(entry),
		Tags:        entry.Tags,
		Note:        entry.Note,
	}, nil
}

// StatusAll returns a KeyStatus for every service in the store,
// sorted alphabetically by service name.
func (s *Store) StatusAll() []KeyStatus {
	names := s.List()
	out := make([]KeyStatus, 0, len(names))
	for _, name := range names {
		if st, err := s.Status(name); err == nil {
			out = append(out, st)
		}
	}
	return out
}
