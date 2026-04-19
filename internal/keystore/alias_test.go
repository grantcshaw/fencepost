package keystore_test

import (
	"testing"

	"github.com/arjunsriva/fencepost/internal/keystore"
)

func TestSetAlias_StoresAlias(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	s.Set("svc", "key123")

	if err := s.SetAlias("svc", "my-service"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	alias, err := s.GetAlias("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if alias != "my-service" {
		t.Errorf("expected %q, got %q", "my-service", alias)
	}
}

func TestSetAlias_MissingService(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	if err := s.SetAlias("ghost", "alias"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestGetAlias_MissingService(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	if _, err := s.GetAlias("ghost"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestClearAlias_RemovesAlias(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	s.Set("svc", "key123")
	s.SetAlias("svc", "my-service")

	if err := s.ClearAlias("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	alias, _ := s.GetAlias("svc")
	if alias != "" {
		t.Errorf("expected empty alias, got %q", alias)
	}
}

func TestServicesByAlias_ReturnsMatchingSorted(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	s.Set("alpha", "k1")
	s.Set("beta", "k2")
	s.Set("gamma", "k3")
	s.SetAlias("alpha", "shared")
	s.SetAlias("gamma", "shared")

	results := s.ServicesByAlias("shared")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0] != "alpha" || results[1] != "gamma" {
		t.Errorf("unexpected order: %v", results)
	}
}

func TestServicesByAlias_NoMatches(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	s.Set("svc", "key")
	results := s.ServicesByAlias("nonexistent")
	if len(results) != 0 {
		t.Errorf("expected no results, got %v", results)
	}
}
