package keystore

import (
	"testing"
)

func TestSetDependencies_StoresAndReturnsSorted(t *testing.T) {
	s := newTestStore(t)
	s.Set("svcA", "keyA")
	s.Set("svcB", "keyB")

	if err := s.SetDependencies("svcA", []string{"svcB", "svcC"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	deps, err := s.GetDependencies("svcA")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(deps) != 2 || deps[0] != "svcB" || deps[1] != "svcC" {
		t.Errorf("unexpected deps: %v", deps)
	}
}

func TestSetDependencies_MissingService(t *testing.T) {
	s := newTestStore(t)
	if err := s.SetDependencies("ghost", []string{"other"}); err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestGetDependencies_MissingService(t *testing.T) {
	s := newTestStore(t)
	if _, err := s.GetDependencies("ghost"); err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestClearDependencies_RemovesDeps(t *testing.T) {
	s := newTestStore(t)
	s.Set("svcA", "keyA")
	s.SetDependencies("svcA", []string{"svcB"})
	if err := s.ClearDependencies("svcA"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	deps, _ := s.GetDependencies("svcA")
	if len(deps) != 0 {
		t.Errorf("expected no deps, got %v", deps)
	}
}

func TestDependentsOf_ReturnsMatchingSorted(t *testing.T) {
	s := newTestStore(t)
	s.Set("svcA", "keyA")
	s.Set("svcB", "keyB")
	s.Set("svcC", "keyC")
	s.SetDependencies("svcB", []string{"svcA"})
	s.SetDependencies("svcC", []string{"svcA"})

	dependents := s.DependentsOf("svcA")
	if len(dependents) != 2 || dependents[0] != "svcB" || dependents[1] != "svcC" {
		t.Errorf("unexpected dependents: %v", dependents)
	}
}

func TestDependentsOf_NoMatches(t *testing.T) {
	s := newTestStore(t)
	s.Set("svcA", "keyA")
	result := s.DependentsOf("svcA")
	if len(result) != 0 {
		t.Errorf("expected empty, got %v", result)
	}
}
