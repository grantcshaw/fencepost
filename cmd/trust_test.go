package cmd_test

import (
	"testing"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func newTrustTestStore(t *testing.T) *keystore.Store {
	t.Helper()
	path := t.TempDir() + "/store.json"
	s := keystore.New(path)
	s.Set("serviceA", "keyA")
	s.Set("serviceB", "keyB")
	return s
}

func TestTrustCmd_SetAndGet(t *testing.T) {
	s := newTrustTestStore(t)

	if err := s.SetTrustLevel("serviceA", "high"); err != nil {
		t.Fatalf("SetTrustLevel failed: %v", err)
	}

	level, err := s.GetTrustLevel("serviceA")
	if err != nil {
		t.Fatalf("GetTrustLevel failed: %v", err)
	}
	if level != "high" {
		t.Errorf("expected high, got %s", level)
	}
}

func TestTrustCmd_DefaultLevel(t *testing.T) {
	s := newTrustTestStore(t)

	level, err := s.GetTrustLevel("serviceB")
	if err != nil {
		t.Fatalf("GetTrustLevel failed: %v", err)
	}
	if level != "none" {
		t.Errorf("expected none as default, got %s", level)
	}
}

func TestTrustCmd_ListByTrustLevel(t *testing.T) {
	s := newTrustTestStore(t)
	s.Set("serviceC", "keyC")
	s.SetTrustLevel("serviceA", "medium")
	s.SetTrustLevel("serviceC", "medium")

	results, err := s.ServicesByTrustLevel("medium")
	if err != nil {
		t.Fatalf("ServicesByTrustLevel failed: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
	if results[0] != "serviceA" || results[1] != "serviceC" {
		t.Errorf("unexpected order: %v", results)
	}
}
