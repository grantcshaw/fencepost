package keystore

import (
	"fmt"
	"sort"
)

type Priority int

const (
	PriorityLow      Priority = 1
	PriorityNormal   Priority = 2
	PriorityHigh     Priority = 3
	PriorityCritical Priority = 4
)

// String returns a human-readable name for the priority level.
func (p Priority) String() string {
	switch p {
	case PriorityLow:
		return "low"
	case PriorityNormal:
		return "normal"
	case PriorityHigh:
		return "high"
	case PriorityCritical:
		return "critical"
	default:
		return fmt.Sprintf("unknown(%d)", int(p))
	}
}

func (s *Store) SetPriority(service string, p Priority) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.data.Keys[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Priority = int(p)
	s.data.Keys[service] = entry
	return s.save()
}

func (s *Store) GetPriority(service string) (Priority, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.data.Keys[service]
	if !ok {
		return 0, fmt.Errorf("service %q not found", service)
	}
	return Priority(entry.Priority), nil
}

func (s *Store) ClearPriority(service string) error {
	return s.SetPriority(service, PriorityNormal)
}

type PriorityEntry struct {
	Service  string
	Priority Priority
}

func (s *Store) ByPriority(p Priority) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var results []string
	for name, entry := range s.data.Keys {
		if Priority(entry.Priority) == p {
			results = append(results, name)
		}
	}
	sort.Strings(results)
	return results
}
