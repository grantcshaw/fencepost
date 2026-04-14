package keystore

import (
	"fmt"
	"time"
)

// RotationPolicy defines when a key should be rotated.
type RotationPolicy struct {
	MaxAgeDays int `json:"max_age_days"`
}

// DefaultRotationPolicy returns a sensible default rotation policy.
func DefaultRotationPolicy() RotationPolicy {
	return RotationPolicy{MaxAgeDays: 90}
}

// DueForRotation returns true if the key for the given service is past its
// rotation deadline according to the provided policy.
func (s *Store) DueForRotation(service string, policy RotationPolicy) (bool, error) {
	entry, err := s.Get(service)
	if err != nil {
		return false, fmt.Errorf("rotation check: %w", err)
	}

	threshold := time.Duration(policy.MaxAgeDays) * 24 * time.Hour
	age := time.Since(entry.RotatedAt)
	return age >= threshold, nil
}

// StaleKeys returns all service names whose keys are due for rotation.
func (s *Store) StaleKeys(policy RotationPolicy) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var stale []string
	threshold := time.Duration(policy.MaxAgeDays) * 24 * time.Hour

	for service, entry := range s.data {
		if time.Since(entry.RotatedAt) >= threshold {
			stale = append(stale, service)
		}
	}
	return stale, nil
}
