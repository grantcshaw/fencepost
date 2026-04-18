package keystore

import "fmt"

// SetFlag marks a service entry with a named flag.
func (s *Store) SetFlag(service, flag string) error {
	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	if entry.Flags == nil {
		entry.Flags = []string{}
	}
	for _, f := range entry.Flags {
		if f == flag {
			return nil
		}
	}
	entry.Flags = append(entry.Flags, flag)
	sortStrings(entry.Flags)
	s.data.Entries[service] = entry
	return s.save()
}

// UnsetFlag removes a named flag from a service entry.
func (s *Store) UnsetFlag(service, flag string) error {
	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	filtered := entry.Flags[:0]
	for _, f := range entry.Flags {
		if f != flag {
			filtered = append(filtered, f)
		}
	}
	entry.Flags = filtered
	s.data.Entries[service] = entry
	return s.save()
}

// GetFlags returns all flags set on a service.
func (s *Store) GetFlags(service string) ([]string, error) {
	entry, ok := s.data.Entries[service]
	if !ok {
		return nil, fmt.Errorf("service %q not found", service)
	}
	return entry.Flags, nil
}

// ServicesByFlag returns all services that have the given flag set.
func (s *Store) ServicesByFlag(flag string) []string {
	var result []string
	for name, entry := range s.data.Entries {
		for _, f := range entry.Flags {
			if f == flag {
				result = append(result, name)
				break
			}
		}
	}
	sortStrings(result)
	return result
}
