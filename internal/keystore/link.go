package keystore

import "fmt"

// SetLink sets a documentation or dashboard URL link for a service.
func (s *Store) SetLink(service, link string) error {
	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Link = link
	s.data.Entries[service] = entry
	return s.save()
}

// GetLink returns the link for a service.
func (s *Store) GetLink(service string) (string, error) {
	entry, ok := s.data.Entries[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}
	return entry.Link, nil
}

// ClearLink removes the link for a service.
func (s *Store) ClearLink(service string) error {
	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Link = ""
	s.data.Entries[service] = entry
	return s.save()
}

// ServicesByLink returns all services that have a non-empty link, sorted.
func (s *Store) ServicesByLink() []string {
	var results []string
	for name, entry := range s.data.Entries {
		if entry.Link != "" {
			results = append(results, name)
		}
	}
	sortStrings(results)
	return results
}
