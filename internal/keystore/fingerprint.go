package keystore

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
)

// SetFingerprint computes and stores a fingerprint for the given service's key.
func (s *Store) SetFingerprint(service string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}

	if entry.Key == "" {
		return "", errors.New("cannot fingerprint empty key")
	}

	hash := sha256.Sum256([]byte(entry.Key))
	fp := hex.EncodeToString(hash[:])
	entry.Fingerprint = fp
	s.data.Entries[service] = entry

	return fp, s.save()
}

// GetFingerprint returns the stored fingerprint for a service.
func (s *Store) GetFingerprint(service string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}
	return entry.Fingerprint, nil
}

// ClearFingerprint removes the fingerprint for a service.
func (s *Store) ClearFingerprint(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}

	entry.Fingerprint = ""
	s.data.Entries[service] = entry
	return s.save()
}

// ServicesByFingerprint returns services whose fingerprint matches the given value.
func (s *Store) ServicesByFingerprint(fp string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []string
	for name, entry := range s.data.Entries {
		if entry.Fingerprint == fp {
			results = append(results, name)
		}
	}
	sort.Strings(results)
	return results
}
