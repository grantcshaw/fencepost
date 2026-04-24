package keystore

import (
	"fmt"
	"sort"
)

// SetMaxRetries sets the maximum retry count for a service's API key operations.
func (s *Store) SetMaxRetries(service string, retries int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	if retries < 0 {
		return fmt.Errorf("retries must be non-negative, got %d", retries)
	}
	entry.MaxRetries = retries
	s.data.Entries[service] = entry
	return s.save()
}

// GetMaxRetries returns the max retries configured for a service.
// Returns 0 and no error if not set.
func (s *Store) GetMaxRetries(service string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return 0, fmt.Errorf("service %q not found", service)
	}
	return entry.MaxRetries, nil
}

// ClearMaxRetries resets the max retries for a service to 0 (unlimited/default).
func (s *Store) ClearMaxRetries(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.MaxRetries = 0
	s.data.Entries[service] = entry
	return s.save()
}

// ServicesByMaxRetries returns all services that have the given max retries value, sorted.
func (s *Store) ServicesByMaxRetries(retries int) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []string
	for name, entry := range s.data.Entries {
		if entry.MaxRetries == retries {
			results = append(results, name)
		}
	}
	sort.Strings(results)
	return results
}
