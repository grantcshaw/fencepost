package keystore

import "fmt"

// CloneResult holds the outcome of a clone operation.
type CloneResult struct {
	Source      string
	Destination string
	Overwritten bool
}

// Clone copies a service entry (key, tags, note) to a new service name.
// If destination already exists, it returns an error unless overwrite is true.
func (s *Store) Clone(src, dst string, overwrite bool) (CloneResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	srcEntry, ok := s.data.Keys[src]
	if !ok {
		return CloneResult{}, fmt.Errorf("service %q not found", src)
	}

	_, exists := s.data.Keys[dst]
	if exists && !overwrite {
		return CloneResult{}, fmt.Errorf("destination service %q already exists; use --overwrite to replace", dst)
	}

	newEntry := Entry{
		Key:       srcEntry.Key,
		CreatedAt: srcEntry.CreatedAt,
		RotatedAt: srcEntry.RotatedAt,
		Note:      srcEntry.Note,
	}

	// Deep-copy tags slice
	if len(srcEntry.Tags) > 0 {
		newEntry.Tags = make([]string, len(srcEntry.Tags))
		copy(newEntry.Tags, srcEntry.Tags)
	}

	s.data.Keys[dst] = newEntry

	if err := s.persist(); err != nil {
		delete(s.data.Keys, dst)
		return CloneResult{}, fmt.Errorf("persist failed: %w", err)
	}

	return CloneResult{
		Source:      src,
		Destination: dst,
		Overwritten: exists,
	}, nil
}
