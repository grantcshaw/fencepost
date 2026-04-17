package keystore

import "fmt"

// SetVersion sets a version string for a service key (e.g. "v2", "2024-01").
func (s *Store) SetVersion(service, version string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Version = version
	s.data.Entries[service] = entry
	return s.save()
}

// GetVersion returns the version string for a service key.
func (s *Store) GetVersion(service string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}
	return entry.Version, nil
}

// ClearVersion removes the version string from a service key.
func (s *Store) ClearVersion(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Version = ""
	s.data.Entries[service] = entry
	return s.save()
}

// ServicesByVersion returns all service names that have the given version, sorted.
func (s *Store) ServicesByVersion(version string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []string
	for name, entry := range s.data.Entries {
		if entry.Version == version {
			results = append(results, name)
		}
	}
	sortResults(results)
	return results
}
