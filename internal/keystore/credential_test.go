package keystore

import (
	"testing"
)

func TestSetCredentialType_StoresType(t *testing.T) {
	s := newTestStore(t)
	s.Set("svc", "key123")
	if err := s.SetCredentialType("svc", "oauth-token"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ct, err := s.GetCredentialType("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ct != "oauth-token" {
		t.Errorf("expected oauth-token, got %s", ct)
	}
}

func TestSetCredentialType_InvalidValue(t *testing.T) {
	s := newTestStore(t)
	s.Set("svc", "key123")
	if err := s.SetCredentialType("svc", "magic-wand"); err == nil {
		t.Error("expected error for invalid credential type")
	}
}

func TestSetCredentialType_MissingService(t *testing.T) {
	s := newTestStore(t)
	if err := s.SetCredentialType("ghost", "jwt"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestGetCredentialType_DefaultsToAPIKey(t *testing.T) {
	s := newTestStore(t)
	s.Set("svc", "key123")
	ct, err := s.GetCredentialType("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ct != "api-key" {
		t.Errorf("expected api-key default, got %s", ct)
	}
}

func TestGetCredentialType_MissingService(t *testing.T) {
	s := newTestStore(t)
	_, err := s.GetCredentialType("ghost")
	if err == nil {
		t.Error("expected error for missing service")
	}
}

func TestServicesByCredentialType_ReturnsMatchingSorted(t *testing.T) {
	s := newTestStore(t)
	for _, svc := range []string{"svc-c", "svc-a", "svc-b"} {
		s.Set(svc, "k")
		s.SetCredentialType(svc, "bearer")
	}
	s.Set("other", "k")
	s.SetCredentialType("other", "jwt")

	results, err := s.ServicesByCredentialType("bearer")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3, got %d", len(results))
	}
	if results[0] != "svc-a" || results[1] != "svc-b" || results[2] != "svc-c" {
		t.Errorf("unexpected order: %v", results)
	}
}
