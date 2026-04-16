package keystore

import "fmt"

// PinEntry marks a service key as pinned, preventing automatic rotation.
func (s *Store) PinEntry(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Pinned = true
	s.data[service] = entry
	return s.save()
}

// UnpinEntry removes the pinned flag from a service key.
func (s *Store) UnpinEntry(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Pinned = false
	s.data[service] = entry
	return s.save()
}

// IsPinned returns true if the service key is pinned.
func (s *Store) IsPinned(service string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data[service]
	if !ok {
		return false, fmt.Errorf("service %q not found", service)
	}
	return entry.Pinned, nil
}

// PinnedKeys returns a sorted list of all pinned service names.
func (s *Store) PinnedKeys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []string
	for name, entry := range s.data {
		if entry.Pinned {
			result = append(result, name)
		}
	}
	sortResults(result)
	return result
}
