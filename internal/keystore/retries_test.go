package keystore_test

import (
	"testing"

	"github.com/cameronbrill/fencepost/internal/keystore"
)

func TestSetMaxRetries_StoresRetries(t *testing.T) {
	s, err := keystore.New(tempStorePath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := s.Set("svc", "key123"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := s.SetMaxRetries("svc", 5); err != nil {
		t.Fatalf("SetMaxRetries: %v", err)
	}
	got, err := s.GetMaxRetries("svc")
	if err != nil {
		t.Fatalf("GetMaxRetries: %v", err)
	}
	if got != 5 {
		t.Errorf("expected 5, got %d", got)
	}
}

func TestSetMaxRetries_MissingService(t *testing.T) {
	s, _ := keystore.New(tempStorePath(t))
	if err := s.SetMaxRetries("ghost", 3); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestSetMaxRetries_NegativeValue(t *testing.T) {
	s, _ := keystore.New(tempStorePath(t))
	_ = s.Set("svc", "key")
	if err := s.SetMaxRetries("svc", -1); err == nil {
		t.Error("expected error for negative retries")
	}
}

func TestGetMaxRetries_MissingService(t *testing.T) {
	s, _ := keystore.New(tempStorePath(t))
	_, err := s.GetMaxRetries("ghost")
	if err == nil {
		t.Error("expected error for missing service")
	}
}

func TestClearMaxRetries_ResetsToZero(t *testing.T) {
	s, _ := keystore.New(tempStorePath(t))
	_ = s.Set("svc", "key")
	_ = s.SetMaxRetries("svc", 7)
	if err := s.ClearMaxRetries("svc"); err != nil {
		t.Fatalf("ClearMaxRetries: %v", err)
	}
	got, _ := s.GetMaxRetries("svc")
	if got != 0 {
		t.Errorf("expected 0 after clear, got %d", got)
	}
}

func TestServicesByMaxRetries_ReturnsMatchingSorted(t *testing.T) {
	s, _ := keystore.New(tempStorePath(t))
	for _, name := range []string{"bravo", "alpha", "charlie"} {
		_ = s.Set(name, "k")
		_ = s.SetMaxRetries(name, 3)
	}
	_ = s.Set("delta", "k")
	_ = s.SetMaxRetries("delta", 10)

	results := s.ServicesByMaxRetries(3)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if results[0] != "alpha" || results[1] != "bravo" || results[2] != "charlie" {
		t.Errorf("unexpected order: %v", results)
	}
}
