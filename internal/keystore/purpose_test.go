package keystore_test

import (
	"testing"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func TestSetPurpose_StoresPurpose(t *testing.T) {
	s := newTestStore(t)
	s.Set("svc", "key123")
	if err := s.SetPurpose("svc", "signing"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := s.GetPurpose("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "signing" {
		t.Errorf("expected signing, got %q", got)
	}
}

func TestSetPurpose_MissingService(t *testing.T) {
	s := newTestStore(t)
	if err := s.SetPurpose("ghost", "signing"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestSetPurpose_InvalidValue(t *testing.T) {
	s := newTestStore(t)
	s.Set("svc", "key123")
	if err := s.SetPurpose("svc", "bogus"); err == nil {
		t.Error("expected error for invalid purpose")
	}
}

func TestGetPurpose_DefaultsToUnknown(t *testing.T) {
	s := newTestStore(t)
	s.Set("svc", "key123")
	got, err := s.GetPurpose("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "unknown" {
		t.Errorf("expected unknown, got %q", got)
	}
}

func TestGetPurpose_MissingService(t *testing.T) {
	s := newTestStore(t)
	_, err := s.GetPurpose("ghost")
	if err == nil {
		t.Error("expected error for missing service")
	}
}

func TestClearPurpose_RemovesPurpose(t *testing.T) {
	s := newTestStore(t)
	s.Set("svc", "key123")
	s.SetPurpose("svc", "encryption")
	if err := s.ClearPurpose("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ := s.GetPurpose("svc")
	if got != "unknown" {
		t.Errorf("expected unknown after clear, got %q", got)
	}
}

func TestServicesByPurpose_ReturnsMatchingSorted(t *testing.T) {
	s := newTestStore(t)
	s.Set("zebra", "k1")
	s.Set("alpha", "k2")
	s.Set("beta", "k3")
	s.SetPurpose("zebra", "webhook")
	s.SetPurpose("alpha", "webhook")
	s.SetPurpose("beta", "signing")
	results, err := s.ServicesByPurpose("webhook")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0] != "alpha" || results[1] != "zebra" {
		t.Errorf("unexpected order: %v", results)
	}
}

func newTestStore(t *testing.T) *keystore.Store {
	t.Helper()
	p := tempStorePath(t)
	s, err := keystore.New(p)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	return s
}
