package keystore_test

import (
	"testing"
	"time"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func TestSchedule_PersistsToDisk(t *testing.T) {
	path := tempStorePath(t)

	s1, err := keystore.New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	_ = s1.Set("svc", "key")
	_ = s1.SetSchedule("svc", "72h")

	s2, err := keystore.New(path)
	if err != nil {
		t.Fatalf("failed to reload store: %v", err)
	}
	got, err := s2.GetSchedule("svc")
	if err != nil {
		t.Fatalf("GetSchedule after reload: %v", err)
	}
	if got != "72h" {
		t.Errorf("expected 72h after reload, got %q", got)
	}
}

func TestSchedule_NextRotation_AfterRotate(t *testing.T) {
	path := tempStorePath(t)
	s, err := keystore.New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	_ = s.Set("svc", "key")
	_ = s.SetSchedule("svc", "24h")

	// Rotate the key
	if err := s.Rotate("svc", "newkey"); err != nil {
		t.Fatalf("Rotate failed: %v", err)
	}

	next, err := s.NextScheduledRotation("svc")
	if err != nil {
		t.Fatalf("NextScheduledRotation failed: %v", err)
	}
	expectedMin := time.Now().Add(23 * time.Hour)
	if next.Before(expectedMin) {
		t.Errorf("expected next rotation at least 23h from now, got %v", next)
	}
}

func TestSchedule_InvalidDuration(t *testing.T) {
	path := tempStorePath(t)
	s, err := keystore.New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	_ = s.Set("svc", "key")
	_ = s.SetSchedule("svc", "not-a-duration")

	_, err = s.NextScheduledRotation("svc")
	if err == nil {
		t.Error("expected error for invalid duration schedule")
	}
}
