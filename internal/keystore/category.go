package keystore

import "fmt"

// SetCategory assigns a category to a service entry.
func (s *Store) SetCategory(service, category string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Category = category
	s.data.Entries[service] = entry
	return s.save()
}

// GetCategory returns the category for a service entry.
func (s *Store) GetCategory(service string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}
	return entry.Category, nil
}

// ClearCategory removes the category from a service entry.
func (s *Store) ClearCategory(service string) error {
	return s.SetCategory(service, "")
}

// ServicesByCategory returns all service names with the given category, sorted.
func (s *Store) ServicesByCategory(category string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []string
	for name, entry := range s.data.Entries {
		if entry.Category == category {
			results = append(results, name)
		}
	}
	sortResults(results)
	return results
}
