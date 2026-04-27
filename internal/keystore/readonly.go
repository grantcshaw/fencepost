package keystore

import "fmt"

// SetReadOnly marks a service's key as read-only, preventing rotation or deletion.
func (s *Store) SetReadOnly(service string, readOnly bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.ReadOnly = readOnly
	s.data.Entries[service] = entry
	return s.save()
}

// IsReadOnly returns true if the service key is marked read-only.
func (s *Store) IsReadOnly(service string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return false, fmt.Errorf("service %q not found", service)
	}
	return entry.ReadOnly, nil
}

// ReadOnlyKeys returns a sorted list of services whose keys are marked read-only.
func (s *Store) ReadOnlyKeys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []string
	for name, entry := range s.data.Entries {
		if entry.ReadOnly {
			results = append(results, name)
		}
	}
	sortStrings(results)
	return results
}
