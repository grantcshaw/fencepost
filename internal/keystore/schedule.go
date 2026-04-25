package keystore

import (
	"fmt"
	"time"
)

// SetSchedule sets the rotation schedule (cron-like string) for a service.
func (s *Store) SetSchedule(service, schedule string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Schedule = schedule
	s.data.Entries[service] = entry
	return s.save()
}

// GetSchedule returns the rotation schedule for a service.
func (s *Store) GetSchedule(service string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}
	return entry.Schedule, nil
}

// ClearSchedule removes the rotation schedule for a service.
func (s *Store) ClearSchedule(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Schedule = ""
	s.data.Entries[service] = entry
	return s.save()
}

// ServicesBySchedule returns all service names that have a non-empty schedule,
// sorted alphabetically.
func (s *Store) ServicesBySchedule() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []string
	for name, entry := range s.data.Entries {
		if entry.Schedule != "" {
			results = append(results, name)
		}
	}
	sortStrings(results)
	return results
}

// NextScheduledRotation parses a simple duration-based schedule (e.g. "24h", "7d")
// and returns the next rotation time based on the key's last rotation.
func (s *Store) NextScheduledRotation(service string) (time.Time, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return time.Time{}, fmt.Errorf("service %q not found", service)
	}
	if entry.Schedule == "" {
		return time.Time{}, fmt.Errorf("service %q has no schedule set", service)
	}
	d, err := time.ParseDuration(entry.Schedule)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid schedule %q: %w", entry.Schedule, err)
	}
	base := entry.RotatedAt
	if base.IsZero() {
		base = entry.CreatedAt
	}
	return base.Add(d), nil
}
