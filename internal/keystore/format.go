package keystore

import "fmt"

// ValidFormats lists all accepted output format values.
var ValidFormats = []string{"json", "yaml", "env", "csv"}

// SetFormat sets the output format preference for a service's key.
func (s *Store) SetFormat(service, format string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !isValidFormat(format) {
		return fmt.Errorf("invalid format %q: must be one of %v", format, ValidFormats)
	}

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}

	entry.Format = format
	s.data.Entries[service] = entry
	return s.save()
}

// GetFormat returns the output format for a service, defaulting to "json".
func (s *Store) GetFormat(service string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}

	if entry.Format == "" {
		return "json", nil
	}
	return entry.Format, nil
}

// ClearFormat resets the format for a service to the default.
func (s *Store) ClearFormat(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}

	entry.Format = ""
	s.data.Entries[service] = entry
	return s.save()
}

// ServicesByFormat returns all service names using the given format, sorted.
func (s *Store) ServicesByFormat(format string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []string
	for name, entry := range s.data.Entries {
		ef := entry.Format
		if ef == "" {
			ef = "json"
		}
		if ef == format {
			results = append(results, name)
		}
	}
	sortResults(results)
	return results
}

func isValidFormat(f string) bool {
	for _, v := range ValidFormats {
		if v == f {
			return true
		}
	}
	return false
}
