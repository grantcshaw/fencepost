package keystore

import "fmt"

// SetOwner assigns an owner string to a service entry.
func (s *Store) SetOwner(service, owner string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Owner = owner
	s.data.Entries[service] = entry
	return s.save()
}

// GetOwner returns the owner assigned to a service entry.
func (s *Store) GetOwner(service string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}
	return entry.Owner, nil
}

// ClearOwner removes the owner from a service entry.
func (s *Store) ClearOwner(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Owner = ""
	s.data.Entries[service] = entry
	return s.save()
}

// ServicesByOwner returns all service names assigned to the given owner, sorted.
func (s *Store) ServicesByOwner(owner string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []string
	for name, entry := range s.data.Entries {
		if entry.Owner == owner {
			results = append(results, name)
		}
	}
	sortResults(results)
	return results
}
