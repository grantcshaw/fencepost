package keystore

import "time"

// LifecycleStatus represents the current lifecycle state of a key.
type LifecycleStatus string

const (
	LifecycleActive     LifecycleStatus = "active"
	LifecycleDeprecated LifecycleStatus = "deprecated"
	LifecycleRetired    LifecycleStatus = "retired"
)

var validLifecycles = map[LifecycleStatus]bool{
	LifecycleActive:     true,
	LifecycleDeprecated: true,
	LifecycleRetired:    true,
}

// SetLifecycle sets the lifecycle status for a service key.
func (s *Store) SetLifecycle(service string, status LifecycleStatus) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.data.Keys[service]
	if !ok {
		return ErrNotFound
	}
	if !validLifecycles[status] {
		return ErrInvalidValue
	}
	entry.Lifecycle = string(status)
	entry.LifecycleUpdatedAt = time.Now().UTC()
	s.data.Keys[service] = entry
	return s.save()
}

// GetLifecycle returns the lifecycle status for a service key.
func (s *Store) GetLifecycle(service string) (LifecycleStatus, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.data.Keys[service]
	if !ok {
		return "", ErrNotFound
	}
	if entry.Lifecycle == "" {
		return LifecycleActive, nil
	}
	return LifecycleStatus(entry.Lifecycle), nil
}

// ServicesByLifecycle returns all services with the given lifecycle status, sorted.
func (s *Store) ServicesByLifecycle(status LifecycleStatus) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var results []string
	for name, entry := range s.data.Keys {
		effective := entry.Lifecycle
		if effective == "" {
			effective = string(LifecycleActive)
		}
		if LifecycleStatus(effective) == status {
			results = append(results, name)
		}
	}
	sortResults(results)
	return results, nil
}
