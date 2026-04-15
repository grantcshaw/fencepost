package keystore

import "fmt"

// Rename copies an entry from oldName to newName and removes the old entry.
// Returns an error if oldName does not exist, or if newName already exists
// and overwrite is false.
func (s *Store) Rename(oldName, newName string, overwrite bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[oldName]
	if !ok {
		return fmt.Errorf("service %q not found", oldName)
	}

	if _, exists := s.data[newName]; exists && !overwrite {
		return fmt.Errorf("service %q already exists; use --overwrite to replace it", newName)
	}

	// Deep-copy tags slice so the two entries don't share the same backing array.
	tagsCopy := make([]string, len(entry.Tags))
	copy(tagsCopy, entry.Tags)

	newEntry := entry
	newEntry.Service = newName
	newEntry.Tags = tagsCopy

	s.data[newName] = newEntry
	delete(s.data, oldName)

	return s.persist()
}
