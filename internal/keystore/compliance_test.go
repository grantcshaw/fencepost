package keystore_test

import (
	"testing"

	"github.com/iamcalledrob/fencepost/internal/keystore"
)

func TestSetCompliance_StoresCompliance(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key123")

	if err := s.SetCompliance("svc", "soc2"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := s.GetCompliance("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "soc2" {
		t.Errorf("expected soc2, got %q", got)
	}
}

func TestSetCompliance_MissingService(t *testing.T) {
	s := newTestStore(t)
	err := s.SetCompliance("ghost", "hipaa")
	if err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestSetCompliance_InvalidValue(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key123")
	err := s.SetCompliance("svc", "made-up-standard")
	if err == nil {
		t.Fatal("expected error for invalid compliance framework")
	}
}

func TestGetCompliance_DefaultsToNone(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key123")

	got, err := s.GetCompliance("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "none" {
		t.Errorf("expected default none, got %q", got)
	}
}

func TestGetCompliance_MissingService(t *testing.T) {
	s := newTestStore(t)
	_, err := s.GetCompliance("ghost")
	if err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestClearCompliance_RemovesCompliance(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key123")
	_ = s.SetCompliance("svc", "gdpr")

	if err := s.ClearCompliance("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ := s.GetCompliance("svc")
	if got != "none" {
		t.Errorf("expected none after clear, got %q", got)
	}
}

func TestServicesByCompliance_ReturnsMatchingSorted(t *testing.T) {
	s := newTestStore(t)
	for _, svc := range []string{"zebra", "alpha", "mango"} {
		_ = s.Set(svc, "k")
		_ = s.SetCompliance(svc, "pci-dss")
	}
	_ = s.Set("other", "k")
	_ = s.SetCompliance("other", "hipaa")

	results := s.ServicesByCompliance("pci-dss")
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if results[0] != "alpha" || results[1] != "mango" || results[2] != "zebra" {
		t.Errorf("unexpected order: %v", results)
	}
}

func TestServicesByCompliance_DefaultNoneIncluded(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("untagged", "k")
	_ = s.Set("explicit-none", "k")
	_ = s.SetCompliance("explicit-none", "none")

	results := s.ServicesByCompliance("none")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func newTestStore(t *testing.T) *keystore.Store {
	t.Helper()
	s, err := keystore.New(tempStorePath(t))
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	return s
}
