package keystore_test

import (
	"testing"
	"time"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func TestSetRevoked_MarksService(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svcA", "key123")

	if err := s.SetRevoked("svcA", "compromised"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ok, err := s.IsRevoked("svcA")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected service to be revoked")
	}
}

func TestSetRevoked_MissingService(t *testing.T) {
	s := newTestStore(t)
	if err := s.SetRevoked("ghost", "reason"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestIsRevoked_DefaultsFalse(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svcA", "key123")

	ok, err := s.IsRevoked("svcA")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected service to not be revoked by default")
	}
}

func TestIsRevoked_MissingService(t *testing.T) {
	s := newTestStore(t)
	_, err := s.IsRevoked("ghost")
	if err == nil {
		t.Error("expected error for missing service")
	}
}

func TestClearRevoked_RemovesRevocation(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svcA", "key123")
	_ = s.SetRevoked("svcA", "test")

	if err := s.ClearRevoked("svcA"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ok, _ := s.IsRevoked("svcA")
	if ok {
		t.Error("expected revoked to be cleared")
	}
}

func TestRevokedKeys_ReturnsSorted(t *testing.T) {
	s := newTestStore(t)
	for _, name := range []string{"gamma", "alpha", "beta"} {
		_ = s.Set(name, "k")
		_ = s.SetRevoked(name, "audit")
	}
	_ = s.Set("delta", "k") // not revoked

	keys := s.RevokedKeys()
	if len(keys) != 3 {
		t.Fatalf("expected 3 revoked keys, got %d", len(keys))
	}
	if keys[0] != "alpha" || keys[1] != "beta" || keys[2] != "gamma" {
		t.Errorf("unexpected order: %v", keys)
	}
}

func TestSetRevoked_PersistsAcrossReload(t *testing.T) {
	path := tempStorePath(t)
	s, _ := keystore.New(path)
	_ = s.Set("svcA", "key123")
	_ = s.SetRevoked("svcA", "leaked")

	s2, _ := keystore.New(path)
	ok, err := s2.IsRevoked("svcA")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected revoked status to persist")
	}
	_ = time.Now() // ensure time import used
}
