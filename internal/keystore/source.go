package keystore

import (
	"fmt"
	"sort"
)

var validSources = map[string]bool{
	"manual":  true,
	"vault":   true,
	"aws":     true,
	"gcp":     true,
	"azure":   true,
	"env":     true,
	"file":    true,
}

func (s *Store) SetSource(service, source string) error {
	if !validSources[source] {
		return fmt.Errorf("invalid source %q: must be one of manual, vault, aws, gcp, azure, env, file", source)
	}
	e, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	e.Source = source
	s.data.Entries[service] = e
	return s.save()
}

func (s *Store) GetSource(service string) (string, error) {
	e, ok := s.data.Entries[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}
	if e.Source == "" {
		return "manual", nil
	}
	return e.Source, nil
}

func (s *Store) ClearSource(service string) error {
	e, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	e.Source = ""
	s.data.Entries[service] = e
	return s.save()
}

func (s *Store) ServicesBySource(source string) ([]string, error) {
	var results []string
	for name, e := range s.data.Entries {
		eff := e.Source
		if eff == "" {
			eff = "manual"
		}
		if eff == source {
			results = append(results, name)
		}
	}
	sort.Strings(results)
	return results, nil
}
