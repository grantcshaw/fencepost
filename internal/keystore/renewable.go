package keystore

import (
	"fmt"
	"sort"
	"time"
)

// SetRenewable marks a service key as auto-renewable or not.
func (s *Store) SetRenewable(service string, renewable bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Renewable = renewable
	s.data.Entries[service] = entry
	return s.save()
}

// IsRenewable returns whether a service key is marked as auto-renewable.
// Defaults to false if not set.
func (s *Store) IsRenewable(service string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return false, fmt.Errorf("service %q not found", service)
	}
	return entry.Renewable, nil
}

// RenewableKeys returns all services marked as auto-renewable, sorted.
func (s *Store) RenewableKeys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []string
	for name, entry := range s.data.Entries {
		if entry.Renewable {
			results = append(results, name)
		}
	}
	sort.Strings(results)
	return results
}

// RenewalDueKeys returns services that are renewable and due for rotation.
func (s *Store) RenewalDueKeys(policy RotationPolicy) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []string
	for name, entry := range s.data.Entries {
		if !entry.Renewable {
			continue
		}
		if policy.DueForRotation(entry.RotatedAt, time.Now()) {
			results = append(results, name)
		}
	}
	sort.Strings(results)
	return results
}
