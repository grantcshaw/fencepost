package keystore_test

import (
	"testing"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func TestSetRating_StoresRating(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	s.Set("svc", "key123")
	if err := s.SetRating("svc", "high"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r, err := s.GetRating("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r != "high" {
		t.Errorf("expected high, got %s", r)
	}
}

func TestSetRating_MissingService(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	if err := s.SetRating("ghost", "low"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestSetRating_InvalidValue(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	s.Set("svc", "key123")
	if err := s.SetRating("svc", "unknown"); err == nil {
		t.Error("expected error for invalid rating")
	}
}

func TestGetRating_DefaultsToMedium(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	s.Set("svc", "key123")
	r, err := s.GetRating("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r != "medium" {
		t.Errorf("expected medium default, got %s", r)
	}
}

func TestGetRating_MissingService(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	if _, err := s.GetRating("ghost"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestServicesByRating_ReturnsMatchingSorted(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	s.Set("alpha", "k1")
	s.Set("beta", "k2")
	s.Set("gamma", "k3")
	s.SetRating("alpha", "critical")
	s.SetRating("gamma", "critical")
	results, err := s.ServicesByRating("critical")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 || results[0] != "alpha" || results[1] != "gamma" {
		t.Errorf("unexpected results: %v", results)
	}
}
