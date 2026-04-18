package keystore_test

import (
	"testing"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func TestSetContact_StoresContact(t *testing.T) {
	ks := keystore.New(tempStorePath(t))
	ks.Set("svc", "key123")

	if err := ks.SetContact("svc", "alice@example.com"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	c, err := ks.GetContact("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c != "alice@example.com" {
		t.Errorf("expected alice@example.com, got %q", c)
	}
}

func TestSetContact_MissingService(t *testing.T) {
	ks := keystore.New(tempStorePath(t))
	if err := ks.SetContact("ghost", "x@y.com"); err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestGetContact_MissingService(t *testing.T) {
	ks := keystore.New(tempStorePath(t))
	_, err := ks.GetContact("ghost")
	if err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestClearContact_RemovesContact(t *testing.T) {
	ks := keystore.New(tempStorePath(t))
	ks.Set("svc", "key123")
	ks.SetContact("svc", "alice@example.com")
	ks.ClearContact("svc")

	c, _ := ks.GetContact("svc")
	if c != "" {
		t.Errorf("expected empty contact, got %q", c)
	}
}

func TestServicesByContact_ReturnsMatchingSorted(t *testing.T) {
	ks := keystore.New(tempStorePath(t))
	ks.Set("zebra", "k1")
	ks.Set("alpha", "k2")
	ks.Set("mango", "k3")
	ks.SetContact("zebra", "team@corp.com")
	ks.SetContact("alpha", "team@corp.com")
	ks.SetContact("mango", "other@corp.com")

	results := ks.ServicesByContact("team@corp.com")
	if len(results) != 2 {
		t.Fatalf("expected 2, got %d", len(results))
	}
	if results[0] != "alpha" || results[1] != "zebra" {
		t.Errorf("unexpected order: %v", results)
	}
}

func TestSetContact_PersistsAcrossReload(t *testing.T) {
	path := tempStorePath(t)
	ks := keystore.New(path)
	ks.Set("svc", "key123")
	ks.SetContact("svc", "persist@test.com")

	ks2 := keystore.New(path)
	c, err := ks2.GetContact("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c != "persist@test.com" {
		t.Errorf("expected persist@test.com, got %q", c)
	}
}
