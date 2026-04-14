package keystore

import (
	"testing"
	"time"
)

func TestDueForRotation_Fresh(t *testing.T) {
	path := tempStorePath(t)
	s, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if err := s.Set("svc", "key123"); err != nil {
		t.Fatalf("Set: %v", err)
	}

	policy := RotationPolicy{MaxAgeDays: 90}
	due, err := s.DueForRotation("svc", policy)
	if err != nil {
		t.Fatalf("DueForRotation: %v", err)
	}
	if due {
		t.Error("expected fresh key to not be due for rotation")
	}
}

func TestDueForRotation_Stale(t *testing.T) {
	path := tempStorePath(t)
	s, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if err := s.Set("svc", "key123"); err != nil {
		t.Fatalf("Set: %v", err)
	}

	// Backdate the RotatedAt timestamp.
	s.mu.Lock()
	e := s.data["svc"]
	e.RotatedAt = time.Now().Add(-100 * 24 * time.Hour)
	s.data["svc"] = e
	s.mu.Unlock()

	policy := RotationPolicy{MaxAgeDays: 90}
	due, err := s.DueForRotation("svc", policy)
	if err != nil {
		t.Fatalf("DueForRotation: %v", err)
	}
	if !due {
		t.Error("expected stale key to be due for rotation")
	}
}

func TestStaleKeys_ReturnsCorrectServices(t *testing.T) {
	path := tempStorePath(t)
	s, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	_ = s.Set("fresh", "k1")
	_ = s.Set("stale", "k2")

	s.mu.Lock()
	e := s.data["stale"]
	e.RotatedAt = time.Now().Add(-100 * 24 * time.Hour)
	s.data["stale"] = e
	s.mu.Unlock()

	policy := RotationPolicy{MaxAgeDays: 90}
	stale, err := s.StaleKeys(policy)
	if err != nil {
		t.Fatalf("StaleKeys: %v", err)
	}
	if len(stale) != 1 || stale[0] != "stale" {
		t.Errorf("expected [stale], got %v", stale)
	}
}
