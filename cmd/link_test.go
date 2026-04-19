package cmd_test

import (
	"testing"

	"github.com/seanmorris/fencepost/internal/keystore"
)

func newLinkTestStore(t *testing.T) (string, *keystore.Store) {
	t.Helper()
	path := t.TempDir() + "/store.json"
	s, err := keystore.New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	return path, s
}

func TestLinkCmd_SetAndGet(t *testing.T) {
	_, s := newLinkTestStore(t)
	if err := s.Set("mysvc", "apikey"); err != nil {
		t.Fatalf("set failed: %v", err)
	}
	if err := s.SetLink("mysvc", "https://docs.example.com"); err != nil {
		t.Fatalf("set link failed: %v", err)
	}
	link, err := s.GetLink("mysvc")
	if err != nil {
		t.Fatalf("get link failed: %v", err)
	}
	if link != "https://docs.example.com" {
		t.Errorf("expected link, got %q", link)
	}
}

func TestLinkCmd_ListByLink(t *testing.T) {
	_, s := newLinkTestStore(t)
	s.Set("svcA", "k1")
	s.Set("svcB", "k2")
	s.Set("svcC", "k3")
	s.SetLink("svcA", "https://a.com")
	s.SetLink("svcC", "https://c.com")
	results := s.ServicesByLink()
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0] != "svcA" || results[1] != "svcC" {
		t.Errorf("unexpected results: %v", results)
	}
}
