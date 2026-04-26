package keystore_test

import (
	"testing"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func TestSetTrustLevel_StoresTrustLevel(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	s.Set("svc", "key123")
	if err := s.SetTrustLevel("svc", "high"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	level, err := s.GetTrustLevel("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if level != "high" {
		t.Errorf("expected high, got %s", level)
	}
}

func TestSetTrustLevel_MissingService(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	err := s.SetTrustLevel("ghost", "low")
	if err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestSetTrustLevel_InvalidValue(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	s.Set("svc", "key123")
	err := s.SetTrustLevel("svc", "extreme")
	if err == nil {
		t.Fatal("expected error for invalid trust level")
	}
}

func TestGetTrustLevel_DefaultsToNone(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	s.Set("svc", "key123")
	level, err := s.GetTrustLevel("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if level != "none" {
		t.Errorf("expected none, got %s", level)
	}
}

func TestGetTrustLevel_MissingService(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	_, err := s.GetTrustLevel("ghost")
	if err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestClearTrustLevel_ResetsToNone(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	s.Set("svc", "key123")
	s.SetTrustLevel("svc", "full")
	if err := s.ClearTrustLevel("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	level, _ := s.GetTrustLevel("svc")
	if level != "none" {
		t.Errorf("expected none after clear, got %s", level)
	}
}

func TestServicesByTrustLevel_ReturnsMatchingSorted(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	s.Set("beta", "k2")
	s.Set("alpha", "k1")
	s.Set("gamma", "k3")
	s.SetTrustLevel("alpha", "high")
	s.SetTrustLevel("beta", "low")
	s.SetTrustLevel("gamma", "high")
	results, err := s.ServicesByTrustLevel("high")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 || results[0] != "alpha" || results[1] != "gamma" {
		t.Errorf("unexpected results: %v", results)
	}
}
