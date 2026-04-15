package keystore

import (
	"fmt"
)

// SetNote attaches a free-text note to a service's key entry.
func (s *Store) SetNote(service, note string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Note = note
	s.data[service] = entry
	return s.save()
}

// GetNote returns the note attached to a service's key entry.
func (s *Store) GetNote(service string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}
	return entry.Note, nil
}

// ClearNote removes the note from a service's key entry.
func (s *Store) ClearNote(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Note = ""
	s.data[service] = entry
	return s.save()
}
