package cmd_test

import (
	"testing"

	"github.com/danielmichaels/fencepost/internal/keystore"
)

func newFlagTestStore(t *testing.T) (string, *keystore.Store) {
	t.Helper()
	path := t.TempDir() + "/store.json"
	ks, err := keystore.New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	return path, ks
}

func TestFlagCmd_SetAndGet(t *testing.T) {
	_, ks := newFlagTestStore(t)
	_ = ks.Set("api", "secret")
	if err := ks.SetFlag("api", "reviewed"); err != nil {
		t.Fatalf("SetFlag error: %v", err)
	}
	flags, err := ks.GetFlags("api")
	if err != nil {
		t.Fatalf("GetFlags error: %v", err)
	}
	if len(flags) != 1 || flags[0] != "reviewed" {
		t.Errorf("expected [reviewed], got %v", flags)
	}
}

func TestFlagCmd_ListByFlag(t *testing.T) {
	_, ks := newFlagTestStore(t)
	_ = ks.Set("svc-x", "k1")
	_ = ks.Set("svc-y", "k2")
	_ = ks.SetFlag("svc-x", "legacy")
	_ = ks.SetFlag("svc-y", "legacy")
	result := ks.ServicesByFlag("legacy")
	if len(result) != 2 {
		t.Errorf("expected 2 services, got %d", len(result))
	}
}
