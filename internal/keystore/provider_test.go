package keystore_test

import (
	"testing"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func TestSetProvider_StoresProvider(t *testing.T) {
	s, err := keystore.New(tempStorePath(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Set("mysvc", "key123"); err != nil {
		t.Fatal(err)
	}
	if err := s.SetProvider("mysvc", "aws"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := s.GetProvider("mysvc")
	if err != nil {
		t.Fatal(err)
	}
	if got != "aws" {
		t.Errorf("expected aws, got %q", got)
	}
}

func TestSetProvider_MissingService(t *testing.T) {
	s, err := keystore.New(tempStorePath(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := s.SetProvider("ghost", "gcp"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestSetProvider_InvalidValue(t *testing.T) {
	s, err := keystore.New(tempStorePath(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Set("mysvc", "key123"); err != nil {
		t.Fatal(err)
	}
	if err := s.SetProvider("mysvc", "notaprovider"); err == nil {
		t.Error("expected error for invalid provider")
	}
}

func TestGetProvider_DefaultsToUnknown(t *testing.T) {
	s, err := keystore.New(tempStorePath(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Set("mysvc", "key123"); err != nil {
		t.Fatal(err)
	}
	got, err := s.GetProvider("mysvc")
	if err != nil {
		t.Fatal(err)
	}
	if got != "unknown" {
		t.Errorf("expected unknown, got %q", got)
	}
}

func TestGetProvider_MissingService(t *testing.T) {
	s, err := keystore.New(tempStorePath(t))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetProvider("ghost"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestServicesByProvider_ReturnsMatchingSorted(t *testing.T) {
	s, err := keystore.New(tempStorePath(t))
	if err != nil {
		t.Fatal(err)
	}
	for _, svc := range []string{"zebra", "alpha", "mango"} {
		if err := s.Set(svc, "k"); err != nil {
			t.Fatal(err)
		}
		if err := s.SetProvider(svc, "stripe"); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Set("other", "k"); err != nil {
		t.Fatal(err)
	}
	if err := s.SetProvider("other", "aws"); err != nil {
		t.Fatal(err)
	}
	results := s.ServicesByProvider("stripe")
	if len(results) != 3 {
		t.Fatalf("expected 3, got %d", len(results))
	}
	if results[0] != "alpha" || results[1] != "mango" || results[2] != "zebra" {
		t.Errorf("unexpected order: %v", results)
	}
}
