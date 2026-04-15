package keystore

import (
	"testing"
	"time"
)

func TestStatus_ReturnsCorrectFields(t *testing.T) {
	s := New(tempStorePath(t))

	if err := s.Set("svc", "key-abc"); err != nil {
		t.Fatalf("Set: %v", err)
	}

	st, err := s.Status("svc")
	if err != nil {
		t.Fatalf("Status: %v", err)
	}

	if st.Service != "svc" {
		t.Errorf("expected service 'svc', got %q", st.Service)
	}
	if st.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
	if st.IsExpired {
		t.Error("fresh key should not be expired")
	}
	if st.DueRotation {
		t.Error("fresh key should not be due for rotation")
	}
}

func TestStatus_MissingService(t *testing.T) {
	s := New(tempStorePath(t))
	_, err := s.Status("ghost")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestStatus_ExpiredKey(t *testing.T) {
	s := New(tempStorePath(t))
	if err := s.Set("old", "key-old"); err != nil {
		t.Fatalf("Set: %v", err)
	}

	// Back-date the entry so it appears expired.
	s.mu.Lock()
	entry := s.data.Keys["old"]
	entry.CreatedAt = time.Now().AddDate(-2, 0, 0)
	s.data.Keys["old"] = entry
	s.mu.Unlock()

	st, err := s.Status("old")
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if !st.IsExpired {
		t.Error("expected key to be expired")
	}
}

func TestStatusAll_ReturnsSortedStatuses(t *testing.T) {
	s := New(tempStorePath(t))
	for _, name := range []string{"zebra", "alpha", "mango"} {
		if err := s.Set(name, "k"); err != nil {
			t.Fatalf("Set %s: %v", name, err)
		}
	}

	all := s.StatusAll()
	if len(all) != 3 {
		t.Fatalf("expected 3 statuses, got %d", len(all))
	}
	expected := []string{"alpha", "mango", "zebra"}
	for i, st := range all {
		if st.Service != expected[i] {
			t.Errorf("index %d: expected %q, got %q", i, expected[i], st.Service)
		}
	}
}
