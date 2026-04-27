package keystore_test

import (
	"testing"

	"github.com/seankim658/fencepost/internal/keystore"
)

func TestRenewable_PersistsToDisk(t *testing.T) {
	path := tempStorePath(t)

	s1, err := keystore.New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	s1.Set("svc", "key")
	if err := s1.SetRenewable("svc", true); err != nil {
		t.Fatalf("SetRenewable failed: %v", err)
	}

	s2, err := keystore.New(path)
	if err != nil {
		t.Fatalf("failed to reload store: %v", err)
	}
	got, err := s2.IsRenewable("svc")
	if err != nil {
		t.Fatalf("IsRenewable failed: %v", err)
	}
	if !got {
		t.Error("expected renewable=true after reload")
	}
}

func TestRenewable_ClearPersistsToDisk(t *testing.T) {
	path := tempStorePath(t)

	s1, err := keystore.New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	s1.Set("svc", "key")
	s1.SetRenewable("svc", true)
	if err := s1.SetRenewable("svc", false); err != nil {
		t.Fatalf("SetRenewable(false) failed: %v", err)
	}

	s2, err := keystore.New(path)
	if err != nil {
		t.Fatalf("failed to reload store: %v", err)
	}
	got, err := s2.IsRenewable("svc")
	if err != nil {
		t.Fatalf("IsRenewable failed: %v", err)
	}
	if got {
		t.Error("expected renewable=false after reload")
	}
}

func TestRenewable_RenewableKeysAfterMultipleSets(t *testing.T) {
	path := tempStorePath(t)

	s, err := keystore.New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	for _, svc := range []string{"a", "b", "c"} {
		s.Set(svc, "k")
		s.SetRenewable(svc, true)
	}
	s.SetRenewable("b", false)

	keys := s.RenewableKeys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 renewable keys, got %d: %v", len(keys), keys)
	}
	if keys[0] != "a" || keys[1] != "c" {
		t.Errorf("unexpected renewable keys: %v", keys)
	}
}
