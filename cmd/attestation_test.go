package cmd

import (
	"testing"

	"github.com/clikd-inc/fencepost/internal/keystore"
)

func newAttestationTestStore(t *testing.T) *keystore.Store {
	t.Helper()
	path := t.TempDir() + "/store.json"
	ks, err := keystore.New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	return ks
}

func TestAttestationCmd_SetAndGet(t *testing.T) {
	ks := newAttestationTestStore(t)
	ks.Set("myapi", "sk-test-123")

	if err := ks.SetAttestation("myapi", "hsm"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	val, err := ks.GetAttestation("myapi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "hsm" {
		t.Errorf("expected hsm, got %s", val)
	}
}

func TestAttestationCmd_DefaultAttestation(t *testing.T) {
	ks := newAttestationTestStore(t)
	ks.Set("myapi", "sk-test-123")

	val, err := ks.GetAttestation("myapi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "none" {
		t.Errorf("expected default none, got %s", val)
	}
}

func TestAttestationCmd_ListByAttestation(t *testing.T) {
	ks := newAttestationTestStore(t)
	ks.Set("svc1", "key-1")
	ks.Set("svc2", "key-2")
	ks.Set("svc3", "key-3")
	ks.SetAttestation("svc1", "tpm")
	ks.SetAttestation("svc3", "tpm")

	results, err := ks.ServicesByAttestation("tpm")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 services, got %d", len(results))
	}
	if results[0] != "svc1" || results[1] != "svc3" {
		t.Errorf("unexpected results: %v", results)
	}
}
