package keystore

import "fmt"

// SetContact sets a contact (owner email/name) for a service key.
func (s *Store) SetContact(service, contact string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Contact = contact
	s.data.Entries[service] = entry
	return s.save()
}

// GetContact returns the contact for a service key.
func (s *Store) GetContact(service string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}
	return entry.Contact, nil
}

// ClearContact removes the contact from a service key.
func (s *Store) ClearContact(service string) error {
	return s.SetContact(service, "")
}

// ServicesByContact returns all service names whose contact matches the given value.
func (s *Store) ServicesByContact(contact string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []string
	for name, entry := range s.data.Entries {
		if entry.Contact == contact {
			results = append(results, name)
		}
	}
	sortStrings(results)
	return results
}
