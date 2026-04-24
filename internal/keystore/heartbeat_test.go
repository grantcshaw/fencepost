package keystore

import (
	"testing"
	"time"
)

func TestSetHeartbeat_StoresTimestamp(t *testing.T) {
	s := newTestStore(t)
	s.Set("svc", "key123")

	now := time.Now().UTC().Truncate(time.Second)
	if err := s.SetHeartbeat("svc", now); err != nil {
		t.Fatalf("SetHeartbeat: %v", err)
	}

	got, err := s.GetHeartbeat("svc")
	if err != nil {
		t.Fatalf("GetHeartbeat: %v", err)
	}
	if !got.Equal(now) {
		t.Errorf("got %v, want %v", got, now)
	}
}

func TestSetHeartbeat_MissingService(t *testing.T) {
	s := newTestStore(t)
	if err := s.SetHeartbeat("ghost", time.Now()); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestGetHeartbeat_MissingService(t *testing.T) {
	s := newTestStore(t)
	if _, err := s.GetHeartbeat("ghost"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestClearHeartbeat_RemovesTimestamp(t *testing.T) {
	s := newTestStore(t)
	s.Set("svc", "key123")
	s.SetHeartbeat("svc", time.Now())

	if err := s.ClearHeartbeat("svc"); err != nil {
		t.Fatalf("ClearHeartbeat: %v", err)
	}

	got, _ := s.GetHeartbeat("svc")
	if !got.IsZero() {
		t.Errorf("expected zero time after clear, got %v", got)
	}
}

func TestSilentServices_ReturnsStaleAndNever(t *testing.T) {
	s := newTestStore(t)
	s.Set("fresh", "k1")
	s.Set("stale", "k2")
	s.Set("never", "k3")

	s.SetHeartbeat("fresh", time.Now())
	s.SetHeartbeat("stale", time.Now().Add(-2*time.Hour))
	// "never" has no heartbeat

	silent := s.SilentServices(30 * time.Minute)

	if len(silent) != 2 {
		t.Fatalf("expected 2 silent services, got %d: %v", len(silent), silent)
	}
	if silent[0] != "never" || silent[1] != "stale" {
		t.Errorf("unexpected silent services: %v", silent)
	}
}

func TestSilentServices_EmptyStore(t *testing.T) {
	s := newTestStore(t)
	result := s.SilentServices(time.Hour)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}
