package keystore_test

import (
	"testing"
	"time"

	"github.com/cameronbroe/fencepost/internal/keystore"
)

func TestSetTimeout_StoresTimeout(t *testing.T) {
	s := newTestStore(t)
	if err := s.Set("svc", "key123"); err != nil {
		t.Fatal(err)
	}
	if err := s.SetTimeout("svc", 30*time.Second); err != nil {
		t.Fatalf("SetTimeout: %v", err)
	}
	got, err := s.GetTimeout("svc")
	if err != nil {
		t.Fatalf("GetTimeout: %v", err)
	}
	if got != 30*time.Second {
		t.Errorf("expected 30s, got %v", got)
	}
}

func TestSetTimeout_MissingService(t *testing.T) {
	s := newTestStore(t)
	if err := s.SetTimeout("ghost", 10*time.Second); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestGetTimeout_MissingService(t *testing.T) {
	s := newTestStore(t)
	if _, err := s.GetTimeout("ghost"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestClearTimeout_RemovesTimeout(t *testing.T) {
	s := newTestStore(t)
	if err := s.Set("svc", "key123"); err != nil {
		t.Fatal(err)
	}
	_ = s.SetTimeout("svc", 5*time.Second)
	if err := s.ClearTimeout("svc"); err != nil {
		t.Fatalf("ClearTimeout: %v", err)
	}
	got, _ := s.GetTimeout("svc")
	if got != 0 {
		t.Errorf("expected 0, got %v", got)
	}
}

func TestServicesByTimeout_ReturnsMatchingSorted(t *testing.T) {
	s := newTestStore(t)
	for _, svc := range []string{"alpha", "beta", "gamma"} {
		if err := s.Set(svc, "k"); err != nil {
			t.Fatal(err)
		}
	}
	_ = s.SetTimeout("alpha", 10*time.Second)
	_ = s.SetTimeout("gamma", 10*time.Second)
	_ = s.SetTimeout("beta", 5*time.Second)

	results := s.ServicesByTimeout(10 * time.Second)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0] != "alpha" || results[1] != "gamma" {
		t.Errorf("unexpected order: %v", results)
	}
}

func TestSetTimeout_PersistsAcrossReload(t *testing.T) {
	path := tempStorePath(t)
	s, err := keystore.New(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Set("svc", "key"); err != nil {
		t.Fatal(err)
	}
	if err := s.SetTimeout("svc", 15*time.Second); err != nil {
		t.Fatal(err)
	}
	s2, err := keystore.New(path)
	if err != nil {
		t.Fatal(err)
	}
	got, err := s2.GetTimeout("svc")
	if err != nil {
		t.Fatal(err)
	}
	if got != 15*time.Second {
		t.Errorf("expected 15s after reload, got %v", got)
	}
}
