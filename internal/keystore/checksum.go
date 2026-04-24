package keystore

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sort"
)

// SetChecksum stores a SHA-256 checksum string for the given service's key.
func (s *Store) SetChecksum(service, checksum string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return errors.New("service not found: " + service)
	}
	entry.Checksum = checksum
	s.data.Entries[service] = entry
	return s.save()
}

// GetChecksum returns the stored checksum for the given service.
func (s *Store) GetChecksum(service string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return "", errors.New("service not found: " + service)
	}
	return entry.Checksum, nil
}

// ClearChecksum removes the stored checksum for the given service.
func (s *Store) ClearChecksum(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return errors.New("service not found: " + service)
	}
	entry.Checksum = ""
	s.data.Entries[service] = entry
	return s.save()
}

// ComputeChecksum calculates a SHA-256 checksum of the service's current key
// value and stores it, returning the hex-encoded digest.
func (s *Store) ComputeChecksum(service string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return "", errors.New("service not found: " + service)
	}
	if entry.Key == "" {
		return "", errors.New("no key set for service: " + service)
	}
	sum := sha256.Sum256([]byte(entry.Key))
	hex := fmt.Sprintf("%x", sum)
	entry.Checksum = hex
	s.data.Entries[service] = entry
	return hex, s.save()
}

// ServicesWithChecksum returns a sorted list of services that have a checksum set.
func (s *Store) ServicesWithChecksum() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []string
	for name, entry := range s.data.Entries {
		if entry.Checksum != "" {
			result = append(result, name)
		}
	}
	sort.Strings(result)
	return result
}
