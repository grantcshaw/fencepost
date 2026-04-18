package keystore

import (
	"errors"
	"sort"
)

var validCredentialTypes = map[string]bool{
	"api-key":      true,
	"oauth-token":  true,
	"service-account": true,
	"bearer":       true,
	"basic":        true,
	"jwt":          true,
}

func (s *Store) SetCredentialType(service, credType string) error {
	if !validCredentialTypes[credType] {
		return errors.New("invalid credential type: " + credType)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.data.Entries[service]
	if !ok {
		return errors.New("service not found: " + service)
	}
	entry.CredentialType = credType
	s.data.Entries[service] = entry
	return s.save()
}

func (s *Store) GetCredentialType(service string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.data.Entries[service]
	if !ok {
		return "", errors.New("service not found: " + service)
	}
	if entry.CredentialType == "" {
		return "api-key", nil
	}
	return entry.CredentialType, nil
}

func (s *Store) ClearCredentialType(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.data.Entries[service]
	if !ok {
		return errors.New("service not found: " + service)
	}
	entry.CredentialType = ""
	s.data.Entries[service] = entry
	return s.save()
}

func (s *Store) ServicesByCredentialType(credType string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var results []string
	for name, entry := range s.data.Entries {
		effective := entry.CredentialType
		if effective == "" {
			effective = "api-key"
		}
		if effective == credType {
			results = append(results, name)
		}
	}
	sort.Strings(results)
	return results, nil
}
