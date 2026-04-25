package keystore_test

import (
	"testing"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func TestSetPolicy_StoresPolicy(t *testing.T) {
	s := newTestStore(t)
	s.Set("svc", "key123")
	if err := s.SetPolicy("svc", "strict"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pol, err := s.GetPolicy("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pol != "strict" {
		t.Errorf("expected strict, got %s", pol)
	}
}

func TestSetPolicy_MissingService(t *testing.T) {
	s := newTestStore(t)
	err := s.SetPolicy("ghost", "strict")
	if err != keystore.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestSetPolicy_InvalidValue(t *testing.T) {
	s := newTestStore(t)
	s.Set("svc", "key123")
	err := s.SetPolicy("svc", "unknown")
	if err == nil {
		t.Error("expected error for invalid policy")
	}
}

func TestGetPolicy_DefaultsToNone(t *testing.T) {
	s := newTestStore(t)
	s.Set("svc", "key123")
	pol, err := s.GetPolicy("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pol != "none" {
		t.Errorf("expected none, got %s", pol)
	}
}

func TestGetPolicy_MissingService(t *testing.T) {
	s := newTestStore(t)
	_, err := s.GetPolicy("ghost")
	if err != keystore.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestClearPolicy_RemovesPolicy(t *testing.T) {
	s := newTestStore(t)
	s.Set("svc", "key123")
	s.SetPolicy("svc", "moderate")
	if err := s.ClearPolicy("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pol, _ := s.GetPolicy("svc")
	if pol != "none" {
		t.Errorf("expected none after clear, got %s", pol)
	}
}

func TestServicesByPolicy_ReturnsMatchingSorted(t *testing.T) {
	s := newTestStore(t)
	s.Set("zebra", "k1")
	s.Set("alpha", "k2")
	s.Set("mango", "k3")
	s.SetPolicy("zebra", "strict")
	s.SetPolicy("alpha", "strict")
	s.SetPolicy("mango", "relaxed")
	result := s.ServicesByPolicy("strict")
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
	if result[0] != "alpha" || result[1] != "zebra" {
		t.Errorf("unexpected order: %v", result)
	}
}
