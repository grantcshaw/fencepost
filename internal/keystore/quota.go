package keystore

import "fmt"

// SetQuota sets the request quota (max calls per period) for a service.
func (s *Store) SetQuota(service string, quota int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Quota = quota
	s.data.Entries[service] = entry
	return s.save()
}

// GetQuota returns the quota for a service. Returns 0 if unset.
func (s *Store) GetQuota(service string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.data.Entries[service]
	if !ok {
		return 0, fmt.Errorf("service %q not found", service)
	}
	return entry.Quota, nil
}

// ClearQuota removes the quota setting for a service.
func (s *Store) ClearQuota(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Quota = 0
	s.data.Entries[service] = entry
	return s.save()
}

// ServicesByQuota returns all services with a quota >= min, sorted by name.
func (s *Store) ServicesByQuota(min int) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []string
	for name, entry := range s.data.Entries {
		if entry.Quota >= min {
			result = append(result, name)
		}
	}
	sortStrings(result)
	return result
}
