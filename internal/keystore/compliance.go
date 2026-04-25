package keystore

import (
	"fmt"
	"sort"
)

var validComplianceFrameworks = map[string]bool{
	"pci-dss":   true,
	"hipaa":     true,
	"soc2":      true,
	"gdpr":      true,
	"iso27001":  true,
	"nist":      true,
	"fedramp":   true,
	"none":      true,
}

// SetCompliance assigns a compliance framework tag to a service's key entry.
func (s *Store) SetCompliance(service, framework string) error {
	if !validComplianceFrameworks[framework] {
		return fmt.Errorf("invalid compliance framework %q: must be one of pci-dss, hipaa, soc2, gdpr, iso27001, nist, fedramp, none", framework)
	}
	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Compliance = framework
	s.data.Entries[service] = entry
	return s.save()
}

// GetCompliance returns the compliance framework for a service.
func (s *Store) GetCompliance(service string) (string, error) {
	entry, ok := s.data.Entries[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}
	if entry.Compliance == "" {
		return "none", nil
	}
	return entry.Compliance, nil
}

// ClearCompliance removes the compliance framework from a service entry.
func (s *Store) ClearCompliance(service string) error {
	entry, ok := s.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Compliance = ""
	s.data.Entries[service] = entry
	return s.save()
}

// ServicesByCompliance returns all services tagged with the given compliance framework, sorted.
func (s *Store) ServicesByCompliance(framework string) []string {
	var results []string
	for name, entry := range s.data.Entries {
		effective := entry.Compliance
		if effective == "" {
			effective = "none"
		}
		if effective == framework {
			results = append(results, name)
		}
	}
	sort.Strings(results)
	return results
}
