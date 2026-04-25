package cmd

import (
	"testing"

	"github.com/janearc/fencepost/internal/keystore"
)

func newRotationCountTestStore(t *testing.T) *keystore.Store {
	t.Helper()
	path := tempStorePath(t)
	s, err := keystore.New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	return s
}

func TestRotationCountGet_DefaultsToZero(t *testing.T) {
	s := newRotationCountTestStore(t)
	_ = s.Set("api", "secret")

	count, err := s.GetRotationCount("api")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0, got %d", count)
	}
}

func TestRotationCountIncrement_ThenGet(t *testing.T) {
	s := newRotationCountTestStore(t)
	_ = s.Set("api", "secret")
	_ = s.IncrementRotationCount("api")
	_ = s.IncrementRotationCount("api")
	_ = s.IncrementRotationCount("api")

	count, err := s.GetRotationCount("api")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 3 {
		t.Errorf("expected 3, got %d", count)
	}
}

func TestRotationCountReset_AfterIncrements(t *testing.T) {
	s := newRotationCountTestStore(t)
	_ = s.Set("api", "secret")
	_ = s.IncrementRotationCount("api")
	_ = s.IncrementRotationCount("api")
	_ = s.ResetRotationCount("api")

	count, _ := s.GetRotationCount("api")
	if count != 0 {
		t.Errorf("expected 0 after reset, got %d", count)
	}
}

func TestRotationCountGet_UnknownKey(t *testing.T) {
	s := newRotationCountTestStore(t)

	_, err := s.GetRotationCount("nonexistent")
	if err == nil {
		t.Error("expected error for unknown key, got nil")
	}
}
