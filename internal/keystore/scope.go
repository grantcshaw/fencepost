package keystore

import (
	"fmt"
	"sort"
)

var validScopes = map[string]bool{
	"read":  true,
	"write": true,
	"admin": true,	
	"read-write": true,
}

func (s *Store) SetScope(service, scope string) error {
	if !validScopes[scope] {
		return fmt.Errorf("invalid scope %q: must be one of read, write, read-write, admin", scope)
	}
	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Scope = scope
	s.data.Entries[service] = entry
	return s.save()
}

func (s *Store) GetScope(service string) (string, error) {
	entry, ok := s.data.Entries[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}
	if entry.Scope == "" {
		return "read", nil
	}
	return entry.Scope, nil
}

func (s *Store) ClearScope(service string) error {
	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Scope = ""
	s.data.Entries[service] = entry
	return s.save()
}

func (s *Store) ServicesByScope(scope string) ([]string, error) {
	var results []string
	for name, entry := range s.data.Entries {
		if entry.Scope == scope {
			results = append(results, name)
		}
	}
	sort.Strings(results)
	return results, nil
}
