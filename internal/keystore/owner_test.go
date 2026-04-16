package keystore

import (
	"testing"
)

func TestSetOwner_StoresOwner(t *testing.T) {
	s, _ := New(tempStorePath(t))
	s.Set("svc", "key123")

	if err := s.SetOwner("svc", "alice"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	owner, err := s.GetOwner("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if owner != "alice" {
		t.Errorf("expected alice, got %q", owner)
	}
}

func TestSetOwner_MissingService(t *testing.T) {
	s, _ := New(tempStorePath(t))
	if err := s.SetOwner("ghost", "alice"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestGetOwner_MissingService(t *testing.T) {
	s, _ := New(tempStorePath(t))
	if _, err := s.GetOwner("ghost"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestClearOwner_RemovesOwner(t *testing.T) {
	s, _ := New(tempStorePath(t))
	s.Set("svc", "key123")
	s.SetOwner("svc", "alice")

	if err := s.ClearOwner("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	owner, _ := s.GetOwner("svc")
	if owner != "" {
		t.Errorf("expected empty owner, got %q", owner)
	}
}

func TestServicesByOwner_ReturnsMatchingSorted(t *testing.T) {
	s, _ := New(tempStorePath(t))
	for _, svc := range []string{"zebra", "alpha", "mango"} {
		s.Set(svc, "k")
		s.SetOwner(svc, "bob")
	}
	s.Set("other", "k")
	s.SetOwner("other", "alice")

	results := s.ServicesByOwner("bob")
	if len(results) != 3 {
		t.Fatalf("expected 3, got %d", len(results))
	}
	if results[0] != "alpha" || results[1] != "mango" || results[2] != "zebra" {
		t.Errorf("unexpected order: %v", results)
	}
}

func TestSetOwner_PersistsAcrossReload(t *testing.T) {
	path := tempStorePath(t)
	s, _ := New(path)
	s.Set("svc", "key")
	s.SetOwner("svc", "carol")

	s2, _ := New(path)
	owner, err := s2.GetOwner("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if owner != "carol" {
		t.Errorf("expected carol, got %q", owner)
	}
}
