package keystore

import (
	"testing"
	"time"
)

func TestHealthCheck_OKKey(t *testing.T) {
	s := New(tempStorePath(t))
	_ = s.Set("svc", "key123")

	r, err := s.HealthCheck("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Status != HealthOK {
		t.Errorf("expected ok, got %s: %v", r.Status, r.Reasons)
	}
}

func TestHealthCheck_MissingService(t *testing.T) {
	s := New(tempStorePath(t))
	_, err := s.HealthCheck("missing")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestHealthCheck_ExpiredKey(t *testing.T) {
	s := New(tempStorePath(t))
	_ = s.Set("svc", "key123")
	entry := s.data.Entries["svc"]
	entry.ExpiresAt = time.Now().Add(-time.Hour)
	s.data.Entries["svc"] = entry

	r, err := s.HealthCheck("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Status != HealthCritical {
		t.Errorf("expected critical, got %s", r.Status)
	}
}

func TestHealthCheck_StaleKey(t *testing.T) {
	s := New(tempStorePath(t))
	_ = s.Set("svc", "key123")
	entry := s.data.Entries["svc"]
	entry.RotatedAt = time.Now().Add(-200 * 24 * time.Hour)
	s.data.Entries["svc"] = entry

	r, err := s.HealthCheck("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Status != HealthWarning && r.Status != HealthCritical {
		t.Errorf("expected warning or critical, got %s", r.Status)
	}
}

func TestHealthCheckAll_ReturnsSortedReports(t *testing.T) {
	s := New(tempStorePath(t))
	_ = s.Set("beta", "k2")
	_ = s.Set("alpha", "k1")

	reports := s.HealthCheckAll()
	if len(reports) != 2 {
		t.Fatalf("expected 2 reports, got %d", len(reports))
	}
	if reports[0].Service != "alpha" || reports[1].Service != "beta" {
		t.Errorf("expected sorted order, got %s %s", reports[0].Service, reports[1].Service)
	}
}
