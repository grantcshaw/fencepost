package cmd_test

import (
	"testing"

	"github.com/bxrne/fencepost/internal/keystore"
)

func newSensitivityTestStore(t *testing.T) *keystore.Store {
	t.Helper()
	s := keystore.New(tempStorePath(t))
	s.Set("api-gateway", "gw-key-abc")
	s.Set("payments", "pay-key-xyz")
	return s
}

func TestSensitivityCmd_SetAndGet(t *testing.T) {
	s := newSensitivityTestStore(t)

	if err := s.SetSensitivity("api-gateway", "confidential"); err != nil {
		t.Fatalf("SetSensitivity failed: %v", err)
	}

	level, err := s.GetSensitivity("api-gateway")
	if err != nil {
		t.Fatalf("GetSensitivity failed: %v", err)
	}
	if level != "confidential" {
		t.Errorf("expected confidential, got %s", level)
	}
}

func TestSensitivityCmd_DefaultLevel(t *testing.T) {
	s := newSensitivityTestStore(t)

	level, err := s.GetSensitivity("payments")
	if err != nil {
		t.Fatalf("GetSensitivity failed: %v", err)
	}
	if level != "internal" {
		t.Errorf("expected default internal, got %s", level)
	}
}

func TestSensitivityCmd_ListBySensitivity(t *testing.T) {
	s := newSensitivityTestStore(t)
	s.Set("auth", "auth-key")

	s.SetSensitivity("api-gateway", "restricted")
	s.SetSensitivity("auth", "restricted")
	s.SetSensitivity("payments", "public")

	results := s.ServicesBySensitivity("restricted")
	if len(results) != 2 {
		t.Fatalf("expected 2 restricted services, got %d", len(results))
	}
	if results[0] != "api-gateway" || results[1] != "auth" {
		t.Errorf("unexpected results: %v", results)
	}
}
