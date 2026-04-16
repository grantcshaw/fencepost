package keystore

import "fmt"

// SetComment sets a short inline comment on a service entry.
func (s *Store) SetComment(service, comment string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Comment = comment
	s.data.Entries[service] = entry
	return s.save()
}

// GetComment returns the comment for a service entry.
func (s *Store) GetComment(service string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}
	return entry.Comment, nil
}

// ClearComment removes the comment from a service entry.
func (s *Store) ClearComment(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Comment = ""
	s.data.Entries[service] = entry
	return s.save()
}
