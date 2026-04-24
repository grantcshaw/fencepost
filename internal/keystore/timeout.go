package keystore

import (
	"errors"
	"sort"
	"time"
)

// SetTimeout sets the request timeout duration for a service's API key.
func (s *Store) SetTimeout(service string, d time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return errors.New("service not found: " + service)
	}
	entry.Timeout = d
	s.data.Entries[service] = entry
	return s.save()
}

// GetTimeout returns the request timeout for a service, or 0 if not set.
func (s *Store) GetTimeout(service string) (time.Duration, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return 0, errors.New("service not found: " + service)
	}
	return entry.Timeout, nil
}

// ClearTimeout removes the timeout setting for a service.
func (s *Store) ClearTimeout(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return errors.New("service not found: " + service)
	}
	entry.Timeout = 0
	s.data.Entries[service] = entry
	return s.save()
}

// ServicesByTimeout returns services whose timeout matches the given duration.
func (s *Store) ServicesByTimeout(d time.Duration) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []string
	for name, entry := range s.data.Entries {
		if entry.Timeout == d {
			results = append(results, name)
		}
	}
	sort.Strings(results)
	return results
}
