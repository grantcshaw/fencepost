package cmd_test

import (
	"testing"

	"github.com/smlrepo/fencepost/internal/keystore"
)

func newTTLTestStore(t *testing.T) (string, *keystore.Store) {
	t.Helper()
	path := t.TempDir() + "/store.json"
	s := keystore.New(path)
	return path, s
}

func TestTTLCmd_SetAndGet(t *testing.T) {
	path, s := newTTLTestStore(t)
	s.Set("mysvc", "apikey")

	if err := s.SetTTL("mysvc", 24); err != nil {
		t.Fatalf("SetTTL failed: %v", err)
	}

	s2 := keystore.New(path)
	hours, err := s2.GetTTL("mysvc")
	if err != nil {
		t.Fatalf("GetTTL failed: %v", err)
	}
	if hours != 24 {
		t.Errorf("expected 24 hours, got %d", hours)
	}
}

func TestTTLCmd_ClearTTL(t *testing.T) {
	path, s := newTTLTestStore(t)
	s.Set("mysvc", "apikey")
	s.SetTTL("mysvc", 48)

	if err := s.ClearTTL("mysvc"); err != nil {
		t.Fatalf("ClearTTL failed: %v", err)
	}

	s2 := keystore.New(path)
	hours, _ := s2.GetTTL("mysvc")
	if hours != 0 {
		t.Errorf("expected 0 after clear, got %d", hours)
	}
}

func TestTTLCmd_ExpiredByTTL_NoneExpired(t *testing.T) {
	_, s := newTTLTestStore(t)
	s.Set("svc1", "k1")
	s.SetTTL("svc1", 9999)

	expired := s.ExpiredByTTL()
	if len(expired) != 0 {
		t.Errorf("expected no expired keys, got %v", expired)
	}
}
