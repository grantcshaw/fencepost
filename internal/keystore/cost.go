package keystore

import (
	"fmt"
	"sort"
)

// SetCost stores a monthly cost estimate (in USD cents) for a service's API key.
func (s *Store) SetCost(service string, cents int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	if cents < 0 {
		return fmt.Errorf("cost must be non-negative")
	}
	entry.Cost = cents
	s.data.Entries[service] = entry
	return s.save()
}

// GetCost returns the monthly cost estimate in USD cents for a service.
func (s *Store) GetCost(service string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return 0, fmt.Errorf("service %q not found", service)
	}
	return entry.Cost, nil
}

// ClearCost resets the cost for a service to zero.
func (s *Store) ClearCost(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Cost = 0
	s.data.Entries[service] = entry
	return s.save()
}

// ServicesByCostAbove returns services whose cost exceeds the given threshold (in cents), sorted by cost descending.
func (s *Store) ServicesByCostAbove(thresholdCents int) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	type pair struct {
		name string
		cost int
	}
	var results []pair
	for name, entry := range s.data.Entries {
		if entry.Cost > thresholdCents {
			results = append(results, pair{name, entry.Cost})
		}
	}
	sort.Slice(results, func(i, j int) bool {
		if results[i].cost != results[j].cost {
			return results[i].cost > results[j].cost
		}
		return results[i].name < results[j].name
	})
	names := make([]string, len(results))
	for i, p := range results {
		names[i] = p.name
	}
	return names
}
