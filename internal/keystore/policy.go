package keystore

import (
	"errors"
	"sort"
)

var validPolicies = map[string]bool{
	"strict":   true,
	"moderate": true,
	"relaxed":  true,
	"none":     true,
}

// SetPolicy assigns a rotation policy to a service.
func (s *Store) SetPolicy(service, policy string) error {
	if !validPolicies[policy] {
		return errors.New("invalid policy: must be one of strict, moderate, relaxed, none")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.data.Entries[service]
	if !ok {
		return ErrNotFound
	}
	entry.Policy = policy
	s.data.Entries[service] = entry
	return s.save()
}

// GetPolicy returns the rotation policy for a service.
func (s *Store) GetPolicy(service string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.data.Entries[service]
	if !ok {
		return "", ErrNotFound
	}
	if entry.Policy == "" {
		return "none", nil
	}
	return entry.Policy, nil
}

// ClearPolicy removes the rotation policy from a service.
func (s *Store) ClearPolicy(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.data.Entries[service]
	if !ok {
		return ErrNotFound
	}
	entry.Policy = ""
	s.data.Entries[service] = entry
	return s.save()
}

// ServicesByPolicy returns all services with the given policy, sorted.
func (s *Store) ServicesByPolicy(policy string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []string
	for name, entry := range s.data.Entries {
		effective := entry.Policy
		if effective == "" {
			effective = "none"
		}
		if effective == policy {
			result = append(result, name)
		}
	}
	sort.Strings(result)
	return result
}
