package keystore

import "fmt"

// CopyKey duplicates a key value from one service to another without removing the source.
// If overwrite is false and destination exists, an error is returned.
func (s *Store) CopyKey(src, dst string, overwrite bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	srcEntry, ok := s.data[src]
	if !ok {
		return fmt.Errorf("source service %q not found", src)
	}

	if _, exists := s.data[dst]; exists && !overwrite {
		return fmt.Errorf("destination service %q already exists; use --overwrite to replace", dst)
	}

	newEntry := Entry{
		Key:       srcEntry.Key,
		CreatedAt: srcEntry.CreatedAt,
		RotatedAt: srcEntry.RotatedAt,
		Tags:      append([]string(nil), srcEntry.Tags...),
		Note:      srcEntry.Note,
	}

	s.data[dst] = newEntry
	return s.save()
}
