package keystore_test

import (
	"testing"
	"time"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func TestSetSchedule_StoresSchedule(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key123")

	if err := s.SetSchedule("svc", "24h"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := s.GetSchedule("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "24h" {
		t.Errorf("expected 24h, got %q", got)
	}
}

func TestSetSchedule_MissingService(t *testing.T) {
	s := newTestStore(t)
	if err := s.SetSchedule("ghost", "24h"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestGetSchedule_MissingService(t *testing.T) {
	s := newTestStore(t)
	_, err := s.GetSchedule("ghost")
	if err == nil {
		t.Error("expected error for missing service")
	}
}

func TestClearSchedule_RemovesSchedule(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key123")
	_ = s.SetSchedule("svc", "48h")

	if err := s.ClearSchedule("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ := s.GetSchedule("svc")
	if got != "" {
		t.Errorf("expected empty schedule, got %q", got)
	}
}

func TestServicesBySchedule_ReturnsMatchingSorted(t *testing.T) {
	s := newTestStore(t)
	for _, svc := range []string{"zebra", "alpha", "mango"} {
		_ = s.Set(svc, "k")
	}
	_ = s.SetSchedule("zebra", "24h")
	_ = s.SetSchedule("alpha", "72h")

	results := s.ServicesBySchedule()
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0] != "alpha" || results[1] != "zebra" {
		t.Errorf("unexpected order: %v", results)
	}
}

func TestNextScheduledRotation_ReturnsCorrectTime(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key123")
	_ = s.SetSchedule("svc", "24h")

	next, err := s.NextScheduledRotation("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if next.IsZero() {
		t.Error("expected non-zero next rotation time")
	}
	if next.Before(time.Now()) {
		t.Error("expected next rotation to be in the future for a fresh key")
	}
}

func TestNextScheduledRotation_NoSchedule(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key123")

	_, err := s.NextScheduledRotation("svc")
	if err == nil {
		t.Error("expected error when no schedule is set")
	}
}

func newTestStore(t *testing.T) *keystore.Store {
	t.Helper()
	path := tempStorePath(t)
	s, err := keystore.New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	return s
}
