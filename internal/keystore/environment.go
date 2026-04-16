package keystore

import "fmt"

// SetEnvironment assigns an environment tag (e.g. "production", "staging") to a service.
func (s *Store) SetEnvironment(service, env string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Environment = env
	s.data.Entries[service] = entry
	return s.save()
}

// GetEnvironment returns the environment assigned to a service.
func (s *Store) GetEnvironment(service string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}
	return entry.Environment, nil
}

// ClearEnvironment removes the environment assignment from a service.
func (s *Store) ClearEnvironment(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Environment = ""
	s.data.Entries[service] = entry
	return s.save()
}

// ServicesByEnvironment returns all service names assigned to the given environment, sorted.
func (s *Store) ServicesByEnvironment(env string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []string
	for name, entry := range s.data.Entries {
		if entry.Environment == env {
			results = append(results, name)
		}
	}
	sortResults(results)
	return results
}
