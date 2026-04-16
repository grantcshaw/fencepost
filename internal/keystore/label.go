package keystore

import "fmt"

// SetLabel assigns a short display label to a service entry.
func (s *Store) SetLabel(service, label string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Label = label
	s.data[service] = entry
	return s.save()
}

// GetLabel returns the label for a service entry.
func (s *Store) GetLabel(service string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}
	return entry.Label, nil
}

// ClearLabel removes the label from a service entry.
func (s *Store) ClearLabel(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Label = ""
	s.data[service] = entry
	return s.save()
}

// LabeledKeys returns all service names that have a non-empty label, sorted.
func (s *Store) LabeledKeys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []string
	for name, entry := range s.data {
		if entry.Label != "" {
			results = append(results, name)
		}
	}
	sortResults(results)
	return results
}
