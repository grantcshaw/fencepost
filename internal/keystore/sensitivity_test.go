package keystore_test

import (
	"testing"

	"github.com/bxrne/fencepost/internal/keystore"
)

func TestSetSensitivity_StoresSensitivity(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	s.Set("svc", "key123")
	if err := s.SetSensitivity("svc", "confidential"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	level, err := s.GetSensitivity("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if level != "confidential" {
		t.Errorf("expected confidential, got %s", level)
	}
}

func TestSetSensitivity_MissingService(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	if err := s.SetSensitivity("ghost", "public"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestSetSensitivity_InvalidValue(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	s.Set("svc", "key123")
	if err := s.SetSensitivity("svc", "top-secret"); err == nil {
		t.Error("expected error for invalid sensitivity level")
	}
}

func TestGetSensitivity_DefaultsToInternal(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	s.Set("svc", "key123")
	level, err := s.GetSensitivity("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if level != "internal" {
		t.Errorf("expected default internal, got %s", level)
	}
}

func TestGetSensitivity_MissingService(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	_, err := s.GetSensitivity("ghost")
	if err == nil {
		t.Error("expected error for missing service")
	}
}

func TestClearSensitivity_RemovesSensitivity(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	s.Set("svc", "key123")
	s.SetSensitivity("svc", "restricted")
	if err := s.ClearSensitivity("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	level, _ := s.GetSensitivity("svc")
	if level != "internal" {
		t.Errorf("expected default internal after clear, got %s", level)
	}
}

func TestServicesBySensitivity_ReturnsMatchingSorted(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	for _, svc := range []string{"alpha", "beta", "gamma"} {
		s.Set(svc, "k")
	}
	s.SetSensitivity("alpha", "restricted")
	s.SetSensitivity("gamma", "restricted")
	s.SetSensitivity("beta", "public")

	results := s.ServicesBySensitivity("restricted")
	if len(results) != 2 {
		t.Fatalf("expected 2, got %d", len(results))
	}
	if results[0] != "alpha" || results[1] != "gamma" {
		t.Errorf("unexpected order: %v", results)
	}
}
