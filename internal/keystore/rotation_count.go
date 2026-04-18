package keystore

import "fmt"

// GetRotationCount returns the number of times a service key has been rotated.
func (s *Store) GetRotationCount(service string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return 0, fmt.Errorf("service %q not found", service)
	}
	return entry.RotationCount, nil
}

// IncrementRotationCount increments the rotation count for a service and persists.
func (s *Store) IncrementRotationCount(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.RotationCount++
	s.data.Entries[service] = entry
	return s.save()
}

// ResetRotationCount resets the rotation count for a service to zero.
func (s *Store) ResetRotationCount(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.RotationCount = 0
	s.data.Entries[service] = entry
	return s.save()
}

// ByRotationCount returns services with at least minCount rotations, sorted by count descending.
func (s *Store) ByRotationCount(minCount int) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	type pair struct {
		name  string
		count int
	}
	var pairs []pair
	for name, entry := range s.data.Entries {
		if entry.RotationCount >= minCount {
			pairs = append(pairs, pair{name, entry.RotationCount})
		}
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count != pairs[j].count {
			return pairs[i].count > pairs[j].count
		}
		return pairs[i].name < pairs[j].name
	})
	result := make([]string, len(pairs))
	for i, p := range pairs {
		result[i] = p.name
	}
	return result
}
