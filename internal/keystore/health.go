package keystore

import "time"

type HealthStatus string

const (
	HealthOK      HealthStatus = "ok"
	HealthWarning HealthStatus = "warning"
	HealthCritical HealthStatus = "critical"
)

type HealthReport struct {
	Service  string
	Status   HealthStatus
	Reasons  []string
}

func (s *Store) HealthCheck(service string) (HealthReport, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data.Entries[service]
	if !ok {
		return HealthReport{}, ErrNotFound
	}

	report := HealthReport{Service: service, Status: HealthOK}

	if entry.Key == "" {
		report.Reasons = append(report.Reasons, "key is empty")
		report.Status = HealthCritical
	}

	if !entry.ExpiresAt.IsZero() && time.Now().After(entry.ExpiresAt) {
		report.Reasons = append(report.Reasons, "key is expired")
		report.Status = HealthCritical
	}

	policy := DefaultRotationPolicy()
	if DueForRotation(entry, policy) {
		report.Reasons = append(report.Reasons, "key is due for rotation")
		if report.Status == HealthOK {
			report.Status = HealthWarning
		}
	}

	if entry.TTL > 0 {
		ttlExpiry := entry.CreatedAt.Add(entry.TTL)
		if time.Now().After(ttlExpiry) {
			report.Reasons = append(report.Reasons, "key TTL has elapsed")
			report.Status = HealthCritical
		}
	}

	return report, nil
}

func (s *Store) HealthCheckAll() []HealthReport {
	s.mu.RLock()
	names := make([]string, 0, len(s.data.Entries))
	for name := range s.data.Entries {
		names = append(names, name)
	}
	s.mu.RUnlock()

	sortStrings(names)
	reports := make([]HealthReport, 0, len(names))
	for _, name := range names {
		if r, err := s.HealthCheck(name); err == nil {
			reports = append(reports, r)
		}
	}
	return reports
}
