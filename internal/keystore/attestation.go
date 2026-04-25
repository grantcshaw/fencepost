package keystore

import (
	"errors"
	"sort"
)

var validAttestations = map[string]bool{
	"none":     true,
	"tpm":      true,
	"hsm":      true,
	"software": true,
	"cloud":    true,
}

// SetAttestation sets the attestation method for a service's key.
func (s *Store) SetAttestation(service, attestation string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !validAttestations[attestation] {
		return errors.New("invalid attestation: must be one of none, tpm, hsm, software, cloud")
	}

	entry, ok := s.data.Entries[service]
	if !ok {
		return errors.New("service not found: " + service)
	}

	entry.Attestation = attestation
	s.data.Entries[service] = entry
	return s.save()
}

// GetAttestation returns the attestation method for a service's key.
func (s *Store) GetAttestation(service string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return "", errors.New("service not found: " + service)
	}

	if entry.Attestation == "" {
		return "none", nil
	}
	return entry.Attestation, nil
}

// ClearAttestation resets the attestation method to none.
func (s *Store) ClearAttestation(service string) error {
	return s.SetAttestation(service, "none")
}

// ServicesByAttestation returns all services using the given attestation method, sorted.
func (s *Store) ServicesByAttestation(attestation string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []string
	for name, entry := range s.data.Entries {
		effective := entry.Attestation
		if effective == "" {
			effective = "none"
		}
		if effective == attestation {
			results = append(results, name)
		}
	}
	sort.Strings(results)
	return results, nil
}
