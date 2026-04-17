package keystore

import (
	"testing"
)

func TestSetLifecycle_StoresLifecycle(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key123")
	if err := s.SetLifecycle("svc", LifecycleDeprecated); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := s.GetLifecycle("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != LifecycleDeprecated {
		t.Errorf("expected deprecated, got %s", got)
	}
}

func TestSetLifecycle_MissingService(t *testing.T) {
	s := newTestStore(t)
	if err := s.SetLifecycle("missing", LifecycleActive); err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestSetLifecycle_InvalidValue(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key123")
	if err := s.SetLifecycle("svc", "unknown"); err != ErrInvalidValue {
		t.Errorf("expected ErrInvalidValue, got %v", err)
	}
}

func TestGetLifecycle_DefaultsToActive(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key123")
	got, err := s.GetLifecycle("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != LifecycleActive {
		t.Errorf("expected active, got %s", got)
	}
}

func TestGetLifecycle_MissingService(t *testing.T) {
	s := newTestStore(t)
	_, err := s.GetLifecycle("missing")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestServicesByLifecycle_ReturnsMatchingSorted(t *testing.T) {
	s := newTestStore(t)
	for _, svc := range []string{"alpha", "beta", "gamma"} {
		_ = s.Set(svc, "k")
	}
	_ = s.SetLifecycle("alpha", LifecycleRetired)
	_ = s.SetLifecycle("gamma", LifecycleRetired)
	results, err := s.ServicesByLifecycle(LifecycleRetired)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 || results[0] != "alpha" || results[1] != "gamma" {
		t.Errorf("unexpected results: %v", results)
	}
}
