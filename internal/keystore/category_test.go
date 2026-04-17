package keystore

import (
	"testing"
)

func TestSetCategory_StoresCategory(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key123")

	if err := s.SetCategory("svc", "payments"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cat, err := s.GetCategory("svc")
	if err != nil || cat != "payments" {
		t.Fatalf("expected payments, got %q %v", cat, err)
	}
}

func TestSetCategory_MissingService(t *testing.T) {
	s := newTestStore(t)
	if err := s.SetCategory("ghost", "infra"); err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestGetCategory_MissingService(t *testing.T) {
	s := newTestStore(t)
	if _, err := s.GetCategory("ghost"); err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestClearCategory_RemovesCategory(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key123")
	_ = s.SetCategory("svc", "payments")
	_ = s.ClearCategory("svc")

	cat, _ := s.GetCategory("svc")
	if cat != "" {
		t.Fatalf("expected empty category, got %q", cat)
	}
}

func TestServicesByCategory_ReturnsMatchingSorted(t *testing.T) {
	s := newTestStore(t)
	for _, svc := range []string{"zebra", "alpha", "mango"} {
		_ = s.Set(svc, "k")
		_ = s.SetCategory(svc, "infra")
	}
	_ = s.Set("other", "k")
	_ = s.SetCategory("other", "payments")

	results := s.ServicesByCategory("infra")
	if len(results) != 3 {
		t.Fatalf("expected 3, got %d", len(results))
	}
	if results[0] != "alpha" || results[1] != "mango" || results[2] != "zebra" {
		t.Fatalf("unexpected order: %v", results)
	}
}

func TestSetCategory_PersistsAcrossReload(t *testing.T) {
	path := tempStorePath(t)
	s, _ := New(path)
	_ = s.Set("svc", "key")
	_ = s.SetCategory("svc", "billing")

	s2, _ := New(path)
	cat, err := s2.GetCategory("svc")
	if err != nil || cat != "billing" {
		t.Fatalf("expected billing after reload, got %q %v", cat, err)
	}
}
