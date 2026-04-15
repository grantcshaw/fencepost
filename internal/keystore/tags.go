package keystore

import "sort"

// SetTags replaces the tags for a given service key entry.
func (s *Store) SetTags(service string, tags []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[service]
	if !ok {
		return ErrNotFound
	}

	normalized := make([]string, len(tags))
	copy(normalized, tags)
	sort.Strings(normalized)

	entry.Tags = normalized
	s.data[service] = entry
	return s.save()
}

// GetTags returns the tags associated with a service key entry.
func (s *Store) GetTags(service string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data[service]
	if !ok {
		return nil, ErrNotFound
	}

	result := make([]string, len(entry.Tags))
	copy(result, entry.Tags)
	return result, nil
}

// FilterByTag returns service names whose entries contain the given tag.
func (s *Store) FilterByTag(tag string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var matches []string
	for service, entry := range s.data {
		for _, t := range entry.Tags {
			if t == tag {
				matches = append(matches, service)
				break
			}
		}
	}
	sort.Strings(matches)
	return matches
}
