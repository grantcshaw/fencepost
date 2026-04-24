package keystore

import (
	"fmt"
	"time"
)

// SetHeartbeat records the last heartbeat timestamp for a service.
func (s *Store) SetHeartbeat(service string, t time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Heartbeat = t
	s.data.Entries[service] = entry
	return s.save()
}

// GetHeartbeat returns the last heartbeat timestamp for a service.
func (s *Store) GetHeartbeat(service string) (time.Time, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return time.Time{}, fmt.Errorf("service %q not found", service)
	}
	return entry.Heartbeat, nil
}

// ClearHeartbeat removes the heartbeat timestamp for a service.
func (s *Store) ClearHeartbeat(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Heartbeat = time.Time{}
	s.data.Entries[service] = entry
	return s.save()
}

// SilentServices returns services whose heartbeat is older than the given
// duration, or that have never had a heartbeat recorded.
func (s *Store) SilentServices(threshold time.Duration) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cutoff := time.Now().Add(-threshold)
	var results []string
	for name, entry := range s.data.Entries {
		if entry.Heartbeat.IsZero() || entry.Heartbeat.Before(cutoff) {
			results = append(results, name)
		}
	}
	sortStrings(results)
	return results
}
