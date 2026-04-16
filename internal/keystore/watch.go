package keystore

import (
	"sort"
	"time"
)

// WatchEntry represents a service being watched for changes.
type WatchEntry struct {
	Service   string    `json:"service"`
	AddedAt   time.Time `json:"added_at"`
}

// Watch adds a service to the watch list.
func (s *Store) Watch(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data.Entries[service]; !ok {
		return ErrNotFound
	}
	if s.data.Watched == nil {
		s.data.Watched = map[string]time.Time{}
	}
	s.data.Watched[service] = time.Now().UTC()
	return s.save()
}

// Unwatch removes a service from the watch list.
func (s *Store) Unwatch(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data.Watched == nil {
		return nil
	}
	delete(s.data.Watched, service)
	return s.save()
}

// IsWatched reports whether a service is being watched.
func (s *Store) IsWatched(service string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.data.Watched == nil {
		return false
	}
	_, ok := s.data.Watched[service]
	return ok
}

// WatchedKeys returns all watched services sorted by name.
func (s *Store) WatchedKeys() []WatchEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var out []WatchEntry
	for svc, t := range s.data.Watched {
		out = append(out, WatchEntry{Service: svc, AddedAt: t})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Service < out[j].Service
	})
	return out
}
