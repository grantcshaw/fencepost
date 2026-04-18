package keystore

import "fmt"

// SetURL sets the base URL associated with a service's API key.
func (s *Store) SetURL(service, url string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.URL = url
	s.data.Entries[service] = entry
	return s.save()
}

// GetURL returns the base URL for a service.
func (s *Store) GetURL(service string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}
	return entry.URL, nil
}

// ClearURL removes the URL from a service entry.
func (s *Store) ClearURL(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.URL = ""
	s.data.Entries[service] = entry
	return s.save()
}

// ServicesByURL returns all service names that have the given URL set.
func (s *Store) ServicesByURL(url string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []string
	for name, entry := range s.data.Entries {
		if entry.URL == url {
			results = append(results, name)
		}
	}
	sortStrings(results)
	return results
}
