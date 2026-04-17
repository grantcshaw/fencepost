package keystore

import "fmt"

// SetRegion assigns a region string to a service entry.
func (s *Store) SetRegion(service, region string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Region = region
	s.data.Entries[service] = entry
	return s.save()
}

// GetRegion returns the region for a service entry.
func (s *Store) GetRegion(service string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}
	return entry.Region, nil
}

// ClearRegion removes the region from a service entry.
func (s *Store) ClearRegion(service string) error {
	return s.SetRegion(service, "")
}

// ServicesByRegion returns all service names that match the given region, sorted.
func (s *Store) ServicesByRegion(region string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []string
	for name, entry := range s.data.Entries {
		if entry.Region == region {
			results = append(results, name)
		}
	}
	sortResults(results)
	return results
}
