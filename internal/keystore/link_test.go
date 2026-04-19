package keystore_test

import (
	"testing"

	"github.com/seanmorris/fencepost/internal/keystore"
)

func TestSetLink_StoresLink(t *testing.T) {
	s, _ := keystore.New(tempStorePath(t))
	s.Set("svc", "key123")
	if err := s.SetLink("svc", "https://example.com/docs"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	link, err := s.GetLink("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if link != "https://example.com/docs" {
		t.Errorf("expected link, got %q", link)
	}
}

func TestSetLink_MissingService(t *testing.T) {
	s, _ := keystore.New(tempStorePath(t))
	if err := s.SetLink("ghost", "https://x.com"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestGetLink_MissingService(t *testing.T) {
	s, _ := keystore.New(tempStorePath(t))
	if _, err := s.GetLink("ghost"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestClearLink_RemovesLink(t *testing.T) {
	s, _ := keystore.New(tempStorePath(t))
	s.Set("svc", "key123")
	s.SetLink("svc", "https://example.com")
	if err := s.ClearLink("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	link, _ := s.GetLink("svc")
	if link != "" {
		t.Errorf("expected empty link, got %q", link)
	}
}

func TestServicesByLink_ReturnsMatchingSorted(t *testing.T) {
	s, _ := keystore.New(tempStorePath(t))
	s.Set("beta", "k2")
	s.Set("alpha", "k1")
	s.Set("gamma", "k3")
	s.SetLink("beta", "https://beta.com")
	s.SetLink("alpha", "https://alpha.com")
	results := s.ServicesByLink()
	if len(results) != 2 {
		t.Fatalf("expected 2, got %d", len(results))
	}
	if results[0] != "alpha" || results[1] != "beta" {
		t.Errorf("unexpected order: %v", results)
	}
}

func TestSetLink_PersistsAcrossReload(t *testing.T) {
	path := tempStorePath(t)
	s, _ := keystore.New(path)
	s.Set("svc", "key123")
	s.SetLink("svc", "https://persist.com")
	s2, _ := keystore.New(path)
	link, err := s2.GetLink("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if link != "https://persist.com" {
		t.Errorf("expected persisted link, got %q", link)
	}
}
