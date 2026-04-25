package keystore

import (
	"fmt"
	"sort"
)

var validProviders = map[string]bool{
	"aws":        true,
	"gcp":        true,
	"azure":      true,
	"github":     true,
	"gitlab":     true,
	"stripe":     true,
	"twilio":     true,
	"sendgrid":   true,
	"datadog":    true,
	"pagerduty":  true,
	"custom":     true,
	"unknown":    true,
}

// SetProvider sets the cloud/service provider for a key entry.
func (s *Store) SetProvider(service, provider string) error {
	if !validProviders[provider] {
		return fmt.Errorf("invalid provider %q: must be one of aws, gcp, azure, github, gitlab, stripe, twilio, sendgrid, datadog, pagerduty, custom, unknown", provider)
	}
	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Provider = provider
	s.data.Entries[service] = entry
	return s.save()
}

// GetProvider returns the provider for a key entry.
func (s *Store) GetProvider(service string) (string, error) {
	entry, ok := s.data.Entries[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}
	if entry.Provider == "" {
		return "unknown", nil
	}
	return entry.Provider, nil
}

// ClearProvider removes the provider field from a key entry.
func (s *Store) ClearProvider(service string) error {
	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Provider = ""
	s.data.Entries[service] = entry
	return s.save()
}

// ServicesByProvider returns all service names with the given provider, sorted.
func (s *Store) ServicesByProvider(provider string) []string {
	var results []string
	for name, entry := range s.data.Entries {
		p := entry.Provider
		if p == "" {
			p = "unknown"
		}
		if p == provider {
			results = append(results, name)
		}
	}
	sort.Strings(results)
	return results
}
