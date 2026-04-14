package keystore

import "sort"

// List returns a sorted slice of all service names in the store.
func (s *Store) List() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	names := make([]string, 0, len(s.data.Keys))
	for name := range s.data.Keys {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
