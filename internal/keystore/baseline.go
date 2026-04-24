package keystore

import (
	"fmt"
	"sort"
	"time"
)

// SetBaseline records the current key value as the baseline for a service.
func (s *Store) SetBaseline(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}

	entry.Baseline = entry.Key
	entry.BaselineAt = time.Now().UTC()
	s.data.Entries[service] = entry
	return s.save()
}

// GetBaseline returns the stored baseline key for a service.
func (s *Store) GetBaseline(service string) (string, time.Time, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return "", time.Time{}, fmt.Errorf("service %q not found", service)
	}

	return entry.Baseline, entry.BaselineAt, nil
}

// ClearBaseline removes the baseline for a service.
func (s *Store) ClearBaseline(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}

	entry.Baseline = ""
	entry.BaselineAt = time.Time{}
	s.data.Entries[service] = entry
	return s.save()
}

// ServicesWithBaseline returns all service names that have a baseline set, sorted.
func (s *Store) ServicesWithBaseline() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var names []string
	for name, entry := range s.data.Entries {
		if entry.Baseline != "" {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names
}

// BaselineChanged reports whether the current key differs from the baseline.
func (s *Store) BaselineChanged(service string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return false, fmt.Errorf("service %q not found", service)
	}

	if entry.Baseline == "" {
		return false, nil
	}
	return entry.Key != entry.Baseline, nil
}
