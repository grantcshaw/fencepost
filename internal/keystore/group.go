package keystore

import (
	"fmt"
	"sort"
)

// SetGroup assigns a named group to a service.
func (s *Store) SetGroup(service, group string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Group = group
	s.data[service] = entry
	return s.save()
}

// GetGroup returns the group assigned to a service.
func (s *Store) GetGroup(service string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}
	return entry.Group, nil
}

// ClearGroup removes the group from a service.
func (s *Store) ClearGroup(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Group = ""
	s.data[service] = entry
	return s.save()
}

// ServicesByGroup returns all services belonging to the given group, sorted.
func (s *Store) ServicesByGroup(group string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []string
	for name, entry := range s.data {
		if entry.Group == group {
			results = append(results, name)
		}
	}
	sort.Strings(results)
	return results
}

// Groups returns a sorted list of all distinct group names in the store.
func (s *Store) Groups() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	seen := make(map[string]struct{})
	for _, entry := range s.data {
		if entry.Group != "" {
			seen[entry.Group] = struct{}{}
		}
	}
	var groups []string
	for g := range seen {
		groups = append(groups, g)
	}
	sort.Strings(groups)
	return groups
}
