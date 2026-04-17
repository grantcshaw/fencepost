package keystore

import (
	"fmt"
	"sort"
	"time"
)

// SetLastAccessed records the current time as the last accessed timestamp for a service.
func (s *Store) SetLastAccessed(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.LastAccessed = time.Now().UTC()
	s.data.Entries[service] = entry
	return s.save()
}

// GetLastAccessed returns the last accessed time for a service.
func (s *Store) GetLastAccessed(service string) (time.Time, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return time.Time{}, fmt.Errorf("service %q not found", service)
	}
	return entry.LastAccessed, nil
}

// NeverAccessed returns a sorted list of services that have never been accessed.
func (s *Store) NeverAccessed() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []string
	for name, entry := range s.data.Entries {
		if entry.LastAccessed.IsZero() {
			result = append(result, name)
		}
	}
	sort.Strings(result)
	return result
}

// AccessedSince returns services accessed after the given time, sorted.
func (s *Store) AccessedSince(t time.Time) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []string
	for name, entry := range s.data.Entries {
		if !entry.LastAccessed.IsZero() && entry.LastAccessed.After(t) {
			result = append(result, name)
		}
	}
	sort.Strings(result)
	return result
}
