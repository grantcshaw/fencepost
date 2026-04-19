package keystore

import "fmt"

// SetAlias assigns a human-friendly alias to a service.
func (s *Store) SetAlias(service, alias string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Alias = alias
	s.data.Entries[service] = entry
	return s.save()
}

// GetAlias returns the alias for a service.
func (s *Store) GetAlias(service string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}
	return entry.Alias, nil
}

// ClearAlias removes the alias from a service.
func (s *Store) ClearAlias(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Alias = ""
	s.data.Entries[service] = entry
	return s.save()
}

// ServicesByAlias returns all services whose alias matches the given value.
func (s *Store) ServicesByAlias(alias string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []string
	for name, entry := range s.data.Entries {
		if entry.Alias == alias {
			results = append(results, name)
		}
	}
	sortStrings(results)
	return results
}
