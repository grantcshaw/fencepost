package keystore

import (
	"testing"
)

func TestSetEnvironment_StoresEnvironment(t *testing.T) {
	s := New(tempStorePath(t))
	s.Set("svcA", "key-1")

	if err := s.SetEnvironment("svcA", "production"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env, err := s.GetEnvironment("svcA")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env != "production" {
		t.Errorf("expected production, got %q", env)
	}
}

func TestSetEnvironment_MissingService(t *testing.T) {
	s := New(tempStorePath(t))
	if err := s.SetEnvironment("ghost", "staging"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestGetEnvironment_MissingService(t *testing.T) {
	s := New(tempStorePath(t))
	if _, err := s.GetEnvironment("ghost"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestClearEnvironment_RemovesEnvironment(t *testing.T) {
	s := New(tempStorePath(t))
	s.Set("svcA", "key-1")
	s.SetEnvironment("svcA", "staging")

	if err := s.ClearEnvironment("svcA"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env, _ := s.GetEnvironment("svcA")
	if env != "" {
		t.Errorf("expected empty environment, got %q", env)
	}
}

func TestServicesByEnvironment_ReturnsMatchingSorted(t *testing.T) {
	s := New(tempStorePath(t))
	s.Set("svcC", "k")
	s.Set("svcA", "k")
	s.Set("svcB", "k")
	s.SetEnvironment("svcC", "production")
	s.SetEnvironment("svcA", "production")
	s.SetEnvironment("svcB", "staging")

	results := s.ServicesByEnvironment("production")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0] != "svcA" || results[1] != "svcC" {
		t.Errorf("unexpected order: %v", results)
	}
}

func TestServicesByEnvironment_PersistsAcrossReload(t *testing.T) {
	path := tempStorePath(t)
	s := New(path)
	s.Set("svcA", "key-1")
	s.SetEnvironment("svcA", "production")

	s2 := New(path)
	env, err := s2.GetEnvironment("svcA")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env != "production" {
		t.Errorf("expected production after reload, got %q", env)
	}
}
