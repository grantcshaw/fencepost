package keystore

import (
	"fmt"
	"sort"
)

// SetDependencies sets the list of services that the given service depends on.
func (s *Store) SetDependencies(service string, deps []string) error {
	if _, ok := s.data.Entries[service]; !ok {
		return fmt.Errorf("service %q not found", service)
	}
	sorted := make([]string, len(deps))
	copy(sorted, deps)
	sort.Strings(sorted)
	e := s.data.Entries[service]
	e.Dependencies = sorted
	s.data.Entries[service] = e
	return s.save()
}

// GetDependencies returns the dependencies for the given service.
func (s *Store) GetDependencies(service string) ([]string, error) {
	e, ok := s.data.Entries[service]
	if !ok {
		return nil, fmt.Errorf("service %q not found", service)
	}
	return e.Dependencies, nil
}

// ClearDependencies removes all dependencies from the given service.
func (s *Store) ClearDependencies(service string) error {
	if _, ok := s.data.Entries[service]; !ok {
		return fmt.Errorf("service %q not found", service)
	}
	e := s.data.Entries[service]
	e.Dependencies = nil
	s.data.Entries[service] = e
	return s.save()
}

// DependentsOf returns all services that list the given service as a dependency.
func (s *Store) DependentsOf(service string) []string {
	var result []string
	for name, e := range s.data.Entries {
		for _, dep := range e.Dependencies {
			if dep == service {
				result = append(result, name)
				break
			}
		}
	}
	sort.Strings(result)
	return result
}
