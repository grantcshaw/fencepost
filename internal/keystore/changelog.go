package keystore

import (
	"fmt"
	"sort"
	"time"
)

// ChangelogEntry records a single change event for a service key.
type ChangelogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     string    `json:"event"`
	Detail    string    `json:"detail,omitempty"`
}

// AppendChangelog adds a new entry to the changelog for the given service.
func (s *Store) AppendChangelog(service, event, detail string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}

	entry.Changelog = append(entry.Changelog, ChangelogEntry{
		Timestamp: time.Now().UTC(),
		Event:     event,
		Detail:    detail,
	})
	s.data.Entries[service] = entry
	return s.save()
}

// GetChangelog returns all changelog entries for a service, oldest first.
func (s *Store) GetChangelog(service string) ([]ChangelogEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return nil, fmt.Errorf("service %q not found", service)
	}
	return entry.Changelog, nil
}

// ClearChangelog removes all changelog entries for a service.
func (s *Store) ClearChangelog(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}

	entry.Changelog = nil
	s.data.Entries[service] = entry
	return s.save()
}

// ServicesWithChangelog returns sorted service names that have at least one changelog entry.
func (s *Store) ServicesWithChangelog() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []string
	for name, entry := range s.data.Entries {
		if len(entry.Changelog) > 0 {
			result = append(result, name)
		}
	}
	sort.Strings(result)
	return result
}
