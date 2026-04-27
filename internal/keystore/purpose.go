package keystore

import "fmt"

var validPurposes = map[string]bool{
	"authentication": true,
	"authorization":  true,
	"encryption":     true,
	"signing":        true,
	"webhook":        true,
	"integration":    true,
	"internal":       true,
	"external":       true,
	"unknown":        true,
}

func (s *Store) SetPurpose(service, purpose string) error {
	if !validPurposes[purpose] {
		return fmt.Errorf("invalid purpose %q: must be one of authentication, authorization, encryption, signing, webhook, integration, internal, external, unknown", purpose)
	}
	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Purpose = purpose
	s.data.Entries[service] = entry
	return s.save()
}

func (s *Store) GetPurpose(service string) (string, error) {
	entry, ok := s.data.Entries[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}
	if entry.Purpose == "" {
		return "unknown", nil
	}
	return entry.Purpose, nil
}

func (s *Store) ClearPurpose(service string) error {
	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Purpose = ""
	s.data.Entries[service] = entry
	return s.save()
}

func (s *Store) ServicesByPurpose(purpose string) ([]string, error) {
	var results []string
	for name, entry := range s.data.Entries {
		effective := entry.Purpose
		if effective == "" {
			effective = "unknown"
		}
		if effective == purpose {
			results = append(results, name)
		}
	}
	sortStrings(results)
	return results, nil
}
