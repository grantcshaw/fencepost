package keystore_test

import (
	"testing"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func TestSetScope_StoresScope(t *testing.T) {
	ks, _ := keystore.New(tempStorePath(t))
	ks.Set("svcA", "key123")

	if err := ks.SetScope("svcA", "admin"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := ks.GetScope("svcA")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "admin" {
		t.Errorf("expected admin, got %s", got)
	}
}

func TestSetScope_MissingService(t *testing.T) {
	ks, _ := keystore.New(tempStorePath(t))
	if err := ks.SetScope("ghost", "read"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestSetScope_InvalidValue(t *testing.T) {
	ks, _ := keystore.New(tempStorePath(t))
	ks.Set("svcA", "key123")
	if err := ks.SetScope("svcA", "superuser"); err == nil {
		t.Error("expected error for invalid scope")
	}
}

func TestGetScope_DefaultsToRead(t *testing.T) {
	ks, _ := keystore.New(tempStorePath(t))
	ks.Set("svcA", "key123")
	got, err := ks.GetScope("svcA")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "read" {
		t.Errorf("expected default read, got %s", got)
	}
}

func TestGetScope_MissingService(t *testing.T) {
	ks, _ := keystore.New(tempStorePath(t))
	if _, err := ks.GetScope("ghost"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestServicesByScope_ReturnsMatchingSorted(t *testing.T) {
	ks, _ := keystore.New(tempStorePath(t))
	for _, svc := range []string{"svcC", "svcA", "svcB"} {
		ks.Set(svc, "key")
		ks.SetScope(svc, "write")
	}
	ks.Set("svcD", "key")
	ks.SetScope("svcD", "admin")

	results, err := ks.ServicesByScope("write")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3, got %d", len(results))
	}
	if results[0] != "svcA" || results[1] != "svcB" || results[2] != "svcC" {
		t.Errorf("unexpected order: %v", results)
	}
}
