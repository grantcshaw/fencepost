package keystore

import (
	"testing"
)

func TestGetRotationCount_DefaultsToZero(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key123")

	count, err := s.GetRotationCount("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0, got %d", count)
	}
}

func TestGetRotationCount_MissingService(t *testing.T) {
	s := newTestStore(t)
	_, err := s.GetRotationCount("ghost")
	if err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestIncrementRotationCount_Increments(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key123")

	_ = s.IncrementRotationCount("svc")
	_ = s.IncrementRotationCount("svc")

	count, _ := s.GetRotationCount("svc")
	if count != 2 {
		t.Errorf("expected 2, got %d", count)
	}
}

func TestResetRotationCount_ResetsToZero(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key123")
	_ = s.IncrementRotationCount("svc")
	_ = s.IncrementRotationCount("svc")
	_ = s.ResetRotationCount("svc")

	count, _ := s.GetRotationCount("svc")
	if count != 0 {
		t.Errorf("expected 0, got %d", count)
	}
}

func TestByRotationCount_ReturnsMatchingSorted(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("alpha", "k1")
	_ = s.Set("beta", "k2")
	_ = s.Set("gamma", "k3")

	_ = s.IncrementRotationCount("alpha")
	_ = s.IncrementRotationCount("alpha")
	_ = s.IncrementRotationCount("beta")

	results := s.ByRotationCount(1)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0] != "alpha" {
		t.Errorf("expected alpha first (highest count), got %s", results[0])
	}
}

func TestIncrementRotationCount_PersistsAcrossReload(t *testing.T) {
	path := tempStorePath(t)
	s1, _ := New(path)
	_ = s1.Set("svc", "key")
	_ = s1.IncrementRotationCount("svc")

	s2, _ := New(path)
	count, _ := s2.GetRotationCount("svc")
	if count != 1 {
		t.Errorf("expected 1 after reload, got %d", count)
	}
}
