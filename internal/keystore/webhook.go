package keystore

import "fmt"

// SetWebhook sets a webhook URL for a service.
func (s *Store) SetWebhook(service, url string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Webhook = url
	s.data.Entries[service] = entry
	return s.save()
}

// GetWebhook returns the webhook URL for a service.
func (s *Store) GetWebhook(service string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}
	return entry.Webhook, nil
}

// ClearWebhook removes the webhook URL for a service.
func (s *Store) ClearWebhook(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Webhook = ""
	s.data.Entries[service] = entry
	return s.save()
}

// ServicesByWebhook returns all services that have a webhook configured, sorted.
func (s *Store) ServicesByWebhook() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []string
	for name, entry := range s.data.Entries {
		if entry.Webhook != "" {
			result = append(result, name)
		}
	}
	sortStrings(result)
	return result
}
