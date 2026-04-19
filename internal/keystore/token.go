package keystore

import "fmt"

// SetToken sets an auth token associated with a service.
func (s *Store) SetToken(service, token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Token = token
	s.data.Entries[service] = entry
	return s.save()
}

// GetToken returns the auth token for a service.
func (s *Store) GetToken(service string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}
	return entry.Token, nil
}

// ClearToken removes the auth token from a service.
func (s *Store) ClearToken(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Token = ""
	s.data.Entries[service] = entry
	return s.save()
}

// ServicesWithToken returns sorted service names that have a token set.
func (s *Store) ServicesWithToken() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []string
	for name, entry := range s.data.Entries {
		if entry.Token != "" {
			result = append(result, name)
		}
	}
	sortStrings(result)
	return result
}
