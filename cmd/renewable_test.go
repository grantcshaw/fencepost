package cmd_test

import (
	"testing"

	"github.com/seankim658/fencepost/internal/keystore"
)

func newRenewableTestStore(t *testing.T) *keystore.Store {
	t.Helper()
	path := tempStorePath(t)
	s, err := keystore.New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	return s
}

func TestRenewableCmd_SetAndGet(t *testing.T) {
	s := newRenewableTestStore(t)
	if err := s.Set("mysvc", "apikey"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	if err := s.SetRenewable("mysvc", true); err != nil {
		t.Fatalf("SetRenewable failed: %v", err)
	}

	got, err := s.IsRenewable("mysvc")
	if err != nil {
		t.Fatalf("IsRenewable failed: %v", err)
	}
	if !got {
		t.Error("expected renewable=true")
	}
}

func TestRenewableCmd_UnsetRenewable(t *testing.T) {
	s := newRenewableTestStore(t)
	s.Set("mysvc", "apikey")
	s.SetRenewable("mysvc", true)
	s.SetRenewable("mysvc", false)

	got, err := s.IsRenewable("mysvc")
	if err != nil {
		t.Fatalf("IsRenewable failed: %v", err)
	}
	if got {
		t.Error("expected renewable=false after unset")
	}
}

func TestRenewableCmd_ListByRenewable(t *testing.T) {
	s := newRenewableTestStore(t)
	for _, svc := range []string{"alpha", "beta", "gamma"} {
		s.Set(svc, "k")
	}
	s.SetRenewable("alpha", true)
	s.SetRenewable("gamma", true)

	keys := s.RenewableKeys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 renewable keys, got %d", len(keys))
	}
	if keys[0] != "alpha" || keys[1] != "gamma" {
		t.Errorf("unexpected keys: %v", keys)
	}
}
