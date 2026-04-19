package keystore

import "fmt"

// SetMetadata stores an arbitrary key-value metadata field for a service.
func (s *Store) SetMetadata(service, key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	if entry.Metadata == nil {
		entry.Metadata = make(map[string]string)
	}
	entry.Metadata[key] = value
	s.data.Entries[service] = entry
	return s.save()
}

// GetMetadata returns the value of a metadata field for a service.
func (s *Store) GetMetadata(service, key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}
	val := entry.Metadata[key]
	return val, nil
}

// ClearMetadata removes a metadata field from a service.
func (s *Store) ClearMetadata(service, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	delete(entry.Metadata, key)
	s.data.Entries[service] = entry
	return s.save()
}

// AllMetadata returns a copy of all metadata fields for a service.
func (s *Store) AllMetadata(service string) (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return nil, fmt.Errorf("service %q not found", service)
	}
	out := make(map[string]string, len(entry.Metadata))
	for k, v := range entry.Metadata {
		out[k] = v
	}
	return out, nil
}
