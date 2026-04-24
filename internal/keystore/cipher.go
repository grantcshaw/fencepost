package keystore

import (
	"errors"
	"fmt"
)

var validCiphers = map[string]bool{
	"aes-256": true,
	"aes-128": true,
	"chacha20": true,
	"none":     true,
}

// SetCipher sets the encryption cipher used for a service's key.
func (s *Store) SetCipher(service, cipher string) error {
	if _, ok := s.data.Entries[service]; !ok {
		return fmt.Errorf("service %q not found", service)
	}
	if !validCiphers[cipher] {
		return fmt.Errorf("invalid cipher %q: must be one of aes-256, aes-128, chacha20, none", cipher)
	}
	e := s.data.Entries[service]
	e.Cipher = cipher
	s.data.Entries[service] = e
	return s.save()
}

// GetCipher returns the cipher for a service, defaulting to "aes-256".
func (s *Store) GetCipher(service string) (string, error) {
	e, ok := s.data.Entries[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}
	if e.Cipher == "" {
		return "aes-256", nil
	}
	return e.Cipher, nil
}

// ClearCipher resets the cipher for a service to the default.
func (s *Store) ClearCipher(service string) error {
	if _, ok := s.data.Entries[service]; !ok {
		return fmt.Errorf("service %q not found", service)
	}
	e := s.data.Entries[service]
	e.Cipher = ""
	s.data.Entries[service] = e
	return s.save()
}

// ServicesByCipher returns all services using the specified cipher, sorted.
func (s *Store) ServicesByCipher(cipher string) ([]string, error) {
	if !validCiphers[cipher] {
		return nil, errors.New("invalid cipher value")
	}
	var result []string
	for name, e := range s.data.Entries {
		eff := e.Cipher
		if eff == "" {
			eff = "aes-256"
		}
		if eff == cipher {
			result = append(result, name)
		}
	}
	sortStrings(result)
	return result, nil
}
