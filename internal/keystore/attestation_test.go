package keystore

import (
	"testing"
)

func TestSetAttestation_StoresAttestation(t *testing.T) {
	s := newTestStore(t)
	s.Set("svcA", "key-abc")

	if err := s.SetAttestation("svcA", "hsm"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	val, err := s.GetAttestation("svcA")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "hsm" {
		t.Errorf("expected hsm, got %s", val)
	}
}

func TestSetAttestation_MissingService(t *testing.T) {
	s := newTestStore(t)

	if err := s.SetAttestation("ghost", "tpm"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestSetAttestation_InvalidValue(t *testing.T) {
	s := newTestStore(t)
	s.Set("svcA", "key-abc")

	if err := s.SetAttestation("svcA", "quantum"); err == nil {
		t.Error("expected error for invalid attestation")
	}
}

func TestGetAttestation_DefaultsToNone(t *testing.T) {
	s := newTestStore(t)
	s.Set("svcA", "key-abc")

	val, err := s.GetAttestation("svcA")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "none" {
		t.Errorf("expected none, got %s", val)
	}
}

func TestGetAttestation_MissingService(t *testing.T) {
	s := newTestStore(t)

	_, err := s.GetAttestation("ghost")
	if err == nil {
		t.Error("expected error for missing service")
	}
}

func TestClearAttestation_ResetsToNone(t *testing.T) {
	s := newTestStore(t)
	s.Set("svcA", "key-abc")
	s.SetAttestation("svcA", "tpm")

	if err := s.ClearAttestation("svcA"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	val, _ := s.GetAttestation("svcA")
	if val != "none" {
		t.Errorf("expected none after clear, got %s", val)
	}
}

func TestServicesByAttestation_ReturnsMatchingSorted(t *testing.T) {
	s := newTestStore(t)
	s.Set("svcA", "key-a")
	s.Set("svcB", "key-b")
	s.Set("svcC", "key-c")
	s.SetAttestation("svcA", "hsm")
	s.SetAttestation("svcC", "hsm")

	results, err := s.ServicesByAttestation("hsm")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0] != "svcA" || results[1] != "svcC" {
		t.Errorf("unexpected order: %v", results)
	}
}
