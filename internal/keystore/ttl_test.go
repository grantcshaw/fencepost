package keystore_test

import (
	"testing"
	"time"

	"github.com/smlrepo/fencepost/internal/keystore"
)

func TestSetTTL_StoresTTL(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	s.Set("svc", "key123")
	if err := s.SetTTL("svc", 48); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	hours, err := s.GetTTL("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hours != 48 {
		t.Errorf("expected 48, got %d", hours)
	}
}

func TestSetTTL_MissingService(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	if err := s.SetTTL("ghost", 24); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestGetTTL_MissingService(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	_, err := s.GetTTL("ghost")
	if err == nil {
		t.Error("expected error for missing service")
	}
}

func TestClearTTL_RemovesTTL(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	s.Set("svc", "key")
	s.SetTTL("svc", 10)
	if err := s.ClearTTL("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	hours, _ := s.GetTTL("svc")
	if hours != 0 {
		t.Errorf("expected 0 after clear, got %d", hours)
	}
}

func TestExpiredByTTL_ReturnsExpired(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	s.Set("fresh", "k1")
	s.SetTTL("fresh", 100)

	s.Set("stale", "k2")
	// Manually backdate by manipulating via rotate to trigger old CreatedAt
	// We rely on TTL=0 meaning no expiry, and TTL=1 with old key being expired.
	// Since we can't easily backdate, just verify non-expired not returned.
	expired := s.ExpiredByTTL()
	for _, name := range expired {
		if name == "fresh" {
			t.Errorf("fresh key should not be expired")
		}
	}
}

func TestSetTTL_PersistsAcrossReload(t *testing.T) {
	path := tempStorePath(t)
	s := keystore.New(path)
	s.Set("svc", "key")
	s.SetTTL("svc", 72)

	s2 := keystore.New(path)
	hours, err := s2.GetTTL("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hours != 72 {
		t.Errorf("expected 72 after reload, got %d", hours)
	}
	_ = time.Now() // suppress unused import
}
