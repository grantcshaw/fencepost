package keystore

import "fmt"

// Valid sensitivity levels
var validSensitivityLevels = map[string]bool{
	"public":       true,
	"internal":     true,
	"confidential": true,
	"restricted":   true,
}

// SetSensitivity assigns a sensitivity level to a service's key entry.
func (s *Store) SetSensitivity(service, level string) error {
	if _, ok := s.data.Entries[service]; !ok {
		return fmt.Errorf("service %q not found", service)
	}
	if !validSensitivityLevels[level] {
		return fmt.Errorf("invalid sensitivity level %q: must be one of public, internal, confidential, restricted", level)
	}
	e := s.data.Entries[service]
	e.Sensitivity = level
	s.data.Entries[service] = e
	return s.save()
}

// GetSensitivity returns the sensitivity level for the given service.
// Defaults to "internal" if not set.
func (s *Store) GetSensitivity(service string) (string, error) {
	e, ok := s.data.Entries[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}
	if e.Sensitivity == "" {
		return "internal", nil
	}
	return e.Sensitivity, nil
}

// ClearSensitivity removes the sensitivity level for the given service,
// reverting it to the default.
func (s *Store) ClearSensitivity(service string) error {
	if _, ok := s.data.Entries[service]; !ok {
		return fmt.Errorf("service %q not found", service)
	}
	e := s.data.Entries[service]
	e.Sensitivity = ""
	s.data.Entries[service] = e
	return s.save()
}

// ServicesBySensitivity returns all service names whose sensitivity matches
// the given level, sorted alphabetically.
func (s *Store) ServicesBySensitivity(level string) []string {
	var results []string
	for name, e := range s.data.Entries {
		effective := e.Sensitivity
		if effective == "" {
			effective = "internal"
		}
		if effective == level {
			results = append(results, name)
		}
	}
	sortStrings(results)
	return results
}
