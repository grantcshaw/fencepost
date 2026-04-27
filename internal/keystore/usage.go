package keystore

import (
	"fmt"
	"sort"
)

// SetUsageCount sets the usage count for a service.
func (s *Store) SetUsageCount(service string, count int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	if count < 0 {
		return fmt.Errorf("usage count must be non-negative")
	}
	entry.UsageCount = count
	s.data.Entries[service] = entry
	return s.save()
}

// IncrementUsageCount increments the usage count for a service by 1.
func (s *Store) IncrementUsageCount(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.UsageCount++
	s.data.Entries[service] = entry
	return s.save()
}

// GetUsageCount returns the usage count for a service.
func (s *Store) GetUsageCount(service string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return 0, fmt.Errorf("service %q not found", service)
	}
	return entry.UsageCount, nil
}

// ResetUsageCount resets the usage count for a service to zero.
func (s *Store) ResetUsageCount(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.UsageCount = 0
	s.data.Entries[service] = entry
	return s.save()
}

// TopByUsage returns service names sorted by usage count descending, up to limit.
// If limit <= 0, all services are returned.
func (s *Store) TopByUsage(limit int) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	type pair struct {
		name  string
		count int
	}
	var pairs []pair
	for name, entry := range s.data.Entries {
		pairs = append(pairs, pair{name, entry.UsageCount})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count != pairs[j].count {
			return pairs[i].count > pairs[j].count
		}
		return pairs[i].name < pairs[j].name
	})
	result := make([]string, 0, len(pairs))
	for _, p := range pairs {
		result = append(result, p.name)
	}
	if limit > 0 && limit < len(result) {
		return result[:limit]
	}
	return result
}
