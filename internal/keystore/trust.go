package keystore

import (
	"fmt"
	"sort"
)

// Valid trust levels
var validTrustLevels = map[string]bool{
	"none":      true,
	"low":       true,
	"medium":    true,
	"high":      true,
	"full":      true,
}

// SetTrustLevel sets the trust level for a service's API key.
func (s *Store) SetTrustLevel(service, level string) error {
	if !validTrustLevels[level] {
		return fmt.Errorf("invalid trust level %q: must be one of none, low, medium, high, full", level)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.TrustLevel = level
	s.data.Entries[service] = entry
	return s.save()
}

// GetTrustLevel returns the trust level for a service, defaulting to "none".
func (s *Store) GetTrustLevel(service string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.data.Entries[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}
	if entry.TrustLevel == "" {
		return "none", nil
	}
	return entry.TrustLevel, nil
}

// ClearTrustLevel resets the trust level for a service to the default ("none").
func (s *Store) ClearTrustLevel(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.TrustLevel = ""
	s.data.Entries[service] = entry
	return s.save()
}

// ServicesByTrustLevel returns all services with the given trust level, sorted.
func (s *Store) ServicesByTrustLevel(level string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var results []string
	for name, entry := range s.data.Entries {
		effective := entry.TrustLevel
		if effective == "" {
			effective = "none"
		}
		if effective == level {
			results = append(results, name)
		}
	}
	sort.Strings(results)
	return results, nil
}
