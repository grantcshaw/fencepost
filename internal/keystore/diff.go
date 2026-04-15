package keystore

import "sort"

// DiffResult describes the difference between two snapshots.
type DiffResult struct {
	Added    []string
	Removed  []string
	Rotated  []string
	Modified []string
}

// Diff compares two stores and returns a DiffResult describing what changed.
// base is the older store; current is the newer one.
func (s *Store) Diff(base map[string]Entry) DiffResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := DiffResult{}

	// Check for added or changed entries
	for name, cur := range s.data {
		base_, exists := base[name]
		if !exists {
			result.Added = append(result.Added, name)
			continue
		}
		if cur.Key != base_.Key {
			if !cur.RotatedAt.IsZero() && cur.RotatedAt.After(base_.RotatedAt) {
				result.Rotated = append(result.Rotated, name)
			} else {
				result.Modified = append(result.Modified, name)
			}
		}
	}

	// Check for removed entries
	for name := range base {
		if _, exists := s.data[name]; !exists {
			result.Removed = append(result.Removed, name)
		}
	}

	sort.Strings(result.Added)
	sort.Strings(result.Removed)
	sort.Strings(result.Rotated)
	sort.Strings(result.Modified)

	return result
}

// SnapshotData returns a shallow copy of the current store entries for diffing.
func (s *Store) SnapshotData() map[string]Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	copy := make(map[string]Entry, len(s.data))
	for k, v := range s.data {
		copy[k] = v
	}
	return copy
}
