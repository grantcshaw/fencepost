package keystore

import (
	"fmt"
	"sort"
	"time"
)

// SetRevoked marks a service key as revoked with an optional reason.
func (s *Store) SetRevoked(service, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}

	entry.Revoked = true
	entry.RevokedAt = time.Now().UTC()
	entry.RevokedReason = reason
	s.data.Entries[service] = entry
	return s.save()
}

// ClearRevoked removes the revoked status from a service key.
func (s *Store) ClearRevoked(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}

	entry.Revoked = false
	entry.RevokedAt = time.Time{}
	entry.RevokedReason = ""
	s.data.Entries[service] = entry
	return s.save()
}

// IsRevoked returns whether a service key has been revoked.
func (s *Store) IsRevoked(service string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return false, fmt.Errorf("service %q not found", service)
	}
	return entry.Revoked, nil
}

// RevokedKeys returns a sorted list of services whose keys are revoked.
func (s *Store) RevokedKeys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []string
	for name, entry := range s.data.Entries {
		if entry.Revoked {
			result = append(result, name)
		}
	}
	sort.Strings(result)
	return result
}
