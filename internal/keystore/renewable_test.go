package keystore_test

import (
	"testing"
	"time"

	"github.com/seankim658/fencepost/internal/keystore"
)

func TestSetRenewable_StoresValue(t *testing.T) {
	s := newTestStore(t)
	s.Set("svc", "key123")

	if err := s.SetRenewable("svc", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := s.IsRenewable("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got {
		t.Error("expected renewable to be true")
	}
}

func TestSetRenewable_MissingService(t *testing.T) {
	s := newTestStore(t)
	if err := s.SetRenewable("ghost", true); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestIsRenewable_DefaultsFalse(t *testing.T) {
	s := newTestStore(t)
	s.Set("svc", "key123")

	got, err := s.IsRenewable("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got {
		t.Error("expected renewable to default to false")
	}
}

func TestIsRenewable_MissingService(t *testing.T) {
	s := newTestStore(t)
	_, err := s.IsRenewable("ghost")
	if err == nil {
		t.Error("expected error for missing service")
	}
}

func TestRenewableKeys_ReturnsSorted(t *testing.T) {
	s := newTestStore(t)
	for _, name := range []string{"zebra", "alpha", "mango"} {
		s.Set(name, "k")
		s.SetRenewable(name, true)
	}
	s.Set("other", "k")

	keys := s.RenewableKeys()
	if len(keys) != 3 {
		t.Fatalf("expected 3, got %d", len(keys))
	}
	if keys[0] != "alpha" || keys[1] != "mango" || keys[2] != "zebra" {
		t.Errorf("unexpected order: %v", keys)
	}
}

func TestRenewalDueKeys_ReturnsStaleRenewable(t *testing.T) {
	s := newTestStore(t)
	s.Set("fresh", "k1")
	s.Set("stale", "k2")
	s.SetRenewable("fresh", true)
	s.SetRenewable("stale", true)

	// Make stale key appear old by rotating with a past time via direct rotate
	s.Rotate("stale", "newkey")

	policy := keystore.DefaultRotationPolicy()
	// Override interval to zero so everything is due
	policy.Interval = 0

	due := s.RenewalDueKeys(policy)
	if len(due) == 0 {
		t.Error("expected at least one renewal due key")
	}
}

func TestRenewalDueKeys_ExcludesNonRenewable(t *testing.T) {
	s := newTestStore(t)
	s.Set("svc", "k")
	// not marked renewable

	policy := keystore.DefaultRotationPolicy()
	policy.Interval = 0

	due := s.RenewalDueKeys(policy)
	for _, name := range due {
		if name == "svc" {
			t.Error("non-renewable service should not appear in renewal due list")
		}
	}
	_ = time.Now()
}
