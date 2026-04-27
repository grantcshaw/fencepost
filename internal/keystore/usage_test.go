package keystore_test

import (
	"testing"

	"github.com/dtnewman/fencepost/internal/keystore"
)

func TestGetUsageCount_DefaultsToZero(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key1")

	count, err := s.GetUsageCount("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0, got %d", count)
	}
}

func TestGetUsageCount_MissingService(t *testing.T) {
	s := newTestStore(t)
	_, err := s.GetUsageCount("ghost")
	if err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestIncrementUsageCount_Increments(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key1")

	for i := 0; i < 3; i++ {
		if err := s.IncrementUsageCount("svc"); err != nil {
			t.Fatalf("increment error: %v", err)
		}
	}
	count, _ := s.GetUsageCount("svc")
	if count != 3 {
		t.Errorf("expected 3, got %d", count)
	}
}

func TestSetUsageCount_NegativeValue(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key1")

	if err := s.SetUsageCount("svc", -1); err == nil {
		t.Fatal("expected error for negative count")
	}
}

func TestResetUsageCount_ResetsToZero(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key1")
	_ = s.SetUsageCount("svc", 10)

	if err := s.ResetUsageCount("svc"); err != nil {
		t.Fatalf("reset error: %v", err)
	}
	count, _ := s.GetUsageCount("svc")
	if count != 0 {
		t.Errorf("expected 0 after reset, got %d", count)
	}
}

func TestTopByUsage_ReturnsSortedDescending(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("alpha", "k1")
	_ = s.Set("beta", "k2")
	_ = s.Set("gamma", "k3")
	_ = s.SetUsageCount("alpha", 5)
	_ = s.SetUsageCount("beta", 20)
	_ = s.SetUsageCount("gamma", 10)

	top := s.TopByUsage(0)
	if len(top) != 3 {
		t.Fatalf("expected 3 results, got %d", len(top))
	}
	if top[0] != "beta" || top[1] != "gamma" || top[2] != "alpha" {
		t.Errorf("unexpected order: %v", top)
	}
}

func TestTopByUsage_RespectsLimit(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("a", "k1")
	_ = s.Set("b", "k2")
	_ = s.Set("c", "k3")
	_ = s.SetUsageCount("a", 1)
	_ = s.SetUsageCount("b", 2)
	_ = s.SetUsageCount("c", 3)

	top := s.TopByUsage(2)
	if len(top) != 2 {
		t.Errorf("expected 2 results, got %d", len(top))
	}
}

func TestUsageCount_PersistsAcrossReload(t *testing.T) {
	path := tempStorePath(t)
	s, _ := keystore.New(path)
	_ = s.Set("svc", "key1")
	_ = s.SetUsageCount("svc", 7)

	s2, _ := keystore.New(path)
	count, err := s2.GetUsageCount("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 7 {
		t.Errorf("expected 7 after reload, got %d", count)
	}
}
