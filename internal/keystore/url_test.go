package keystore_test

import (
	"testing"

	"github.com/clikd-inc/fencepost/internal/keystore"
)

func TestSetURL_StoresURL(t *testing.T) {
	ks := keystore.New(tempStorePath(t))
	ks.Set("mysvc", "key123")

	if err := ks.SetURL("mysvc", "https://api.example.com"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := ks.GetURL("mysvc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "https://api.example.com" {
		t.Errorf("expected %q, got %q", "https://api.example.com", got)
	}
}

func TestSetURL_MissingService(t *testing.T) {
	ks := keystore.New(tempStorePath(t))
	if err := ks.SetURL("ghost", "https://x.com"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestGetURL_MissingService(t *testing.T) {
	ks := keystore.New(tempStorePath(t))
	if _, err := ks.GetURL("ghost"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestClearURL_RemovesURL(t *testing.T) {
	ks := keystore.New(tempStorePath(t))
	ks.Set("mysvc", "key123")
	ks.SetURL("mysvc", "https://api.example.com")

	if err := ks.ClearURL("mysvc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ := ks.GetURL("mysvc")
	if got != "" {
		t.Errorf("expected empty URL, got %q", got)
	}
}

func TestServicesByURL_ReturnsMatchingSorted(t *testing.T) {
	ks := keystore.New(tempStorePath(t))
	for _, svc := range []string{"svcC", "svcA", "svcB"} {
		ks.Set(svc, "k")
		ks.SetURL(svc, "https://shared.api")
	}
	ks.Set("svcD", "k")
	ks.SetURL("svcD", "https://other.api")

	results := ks.ServicesByURL("https://shared.api")
	if len(results) != 3 {
		t.Fatalf("expected 3, got %d", len(results))
	}
	if results[0] != "svcA" || results[1] != "svcB" || results[2] != "svcC" {
		t.Errorf("unexpected order: %v", results)
	}
}
