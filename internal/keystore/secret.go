package keystore

import "fmt"

// SetSecret stores an arbitrary secondary secret (e.g. client secret) for a service.
func (s *Store) SetSecret(service, secret string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Secret = secret
	s.data.Entries[service] = entry
	return s.save()
}

// GetSecret returns the secondary secret for a service.
func (s *Store) GetSecret(service string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}
	return entry.Secret, nil
}

// ClearSecret removes the secondary secret for a service.
func (s *Store) ClearSecret(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Secret = ""
	s.data.Entries[service] = entry
	return s.save()
}

// ServicesWithSecret returns all service names that have a non-empty secret, sorted.
func (s *Store) ServicesWithSecret() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var names []string
	for name, entry := range s.data.Entries {
		if entry.Secret != "" {
			names = append(names, name)
		}
	}
	sortStrings(names)
	return names
}
