package keystore

import (
	"fmt"
	"time"
)

// SetTTL sets a time-to-live duration (in hours) for a service key.
func (s *Store) SetTTL(service string, hours int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	if hours < 0 {
		return fmt.Errorf("TTL must be non-negative")
	}
	entry.TTLHours = hours
	s.data.Entries[service] = entry
	return s.save()
}

// GetTTL returns the TTL in hours for a service key (0 means no TTL set).
func (s *Store) GetTTL(service string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return 0, fmt.Errorf("service %q not found", service)
	}
	return entry.TTLHours, nil
}

// ClearTTL removes the TTL for a service key.
func (s *Store) ClearTTL(service string) error {
	return s.SetTTL(service, 0)
}

// ExpiredByTTL returns services whose key has exceeded its TTL since creation.
func (s *Store) ExpiredByTTL() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := time.Now()
	var result []string
	for name, entry := range s.data.Entries {
		if entry.TTLHours > 0 && !entry.CreatedAt.IsZero() {
			deadline := entry.CreatedAt.Add(time.Duration(entry.TTLHours) * time.Hour)
			if now.After(deadline) {
				result = append(result, name)
			}
		}
	}
	sort.Strings(result)
	return result
}
