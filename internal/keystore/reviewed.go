package keystore

import (
	"fmt"
	"sort"
	"time"
)

// SetReviewedAt records the last manual review timestamp for a service.
func (s *Store) SetReviewedAt(service string, t time.Time) error {
	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.ReviewedAt = t
	s.data.Entries[service] = entry
	return s.save()
}

// GetReviewedAt returns the last review timestamp for a service.
func (s *Store) GetReviewedAt(service string) (time.Time, error) {
	entry, ok := s.data.Entries[service]
	if !ok {
		return time.Time{}, fmt.Errorf("service %q not found", service)
	}
	return entry.ReviewedAt, nil
}

// ClearReviewedAt removes the review timestamp for a service.
func (s *Store) ClearReviewedAt(service string) error {
	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.ReviewedAt = time.Time{}
	s.data.Entries[service] = entry
	return s.save()
}

// NeverReviewed returns service names that have no review timestamp.
func (s *Store) NeverReviewed() []string {
	var result []string
	for name, entry := range s.data.Entries {
		if entry.ReviewedAt.IsZero() {
			result = append(result, name)
		}
	}
	sort.Strings(result)
	return result
}

// ReviewedBefore returns services whose last review is older than the given time.
func (s *Store) ReviewedBefore(cutoff time.Time) []string {
	var result []string
	for name, entry := range s.data.Entries {
		if !entry.ReviewedAt.IsZero() && entry.ReviewedAt.Before(cutoff) {
			result = append(result, name)
		}
	}
	sort.Strings(result)
	return result
}
