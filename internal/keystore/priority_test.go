package keystore_test

import (
	"testing"

	"github.com/nqzyx/fencepost/internal/keystore"
)

func TestSetPriority_StoresPriority(t *testing.T) {
	s, _ := keystore.New(tempStorePath(t))
	s.Set("svc", "key123")
	if err := s.SetPriority("svc", keystore.PriorityHigh); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p, err := s.GetPriority("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p != keystore.PriorityHigh {
		t.Errorf("expected PriorityHigh, got %v", p)
	}
}

func TestSetPriority_MissingService(t *testing.T) {
	s, _ := keystore.New(tempStorePath(t))
	if err := s.SetPriority("ghost", keystore.PriorityLow); err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestGetPriority_MissingService(t *testing.T) {
	s, _ := keystore.New(tempStorePath(t))
	if _, err := s.GetPriority("ghost"); err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestClearPriority_ResetsToNormal(t *testing.T) {
	s, _ := keystore.New(tempStorePath(t))
	s.Set("svc", "key123")
	s.SetPriority("svc", keystore.PriorityCritical)
	s.ClearPriority("svc")
	p, _ := s.GetPriority("svc")
	if p != keystore.PriorityNormal {
		t.Errorf("expected PriorityNormal after clear, got %v", p)
	}
}

func TestByPriority_ReturnsMatchingSorted(t *testing.T) {
	s, _ := keystore.New(tempStorePath(t))
	s.Set("alpha", "k1")
	s.Set("beta", "k2")
	s.Set("gamma", "k3")
	s.SetPriority("alpha", keystore.PriorityHigh)
	s.SetPriority("gamma", keystore.PriorityHigh)
	results := s.ByPriority(keystore.PriorityHigh)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0] != "alpha" || results[1] != "gamma" {
		t.Errorf("unexpected order: %v", results)
	}
}
