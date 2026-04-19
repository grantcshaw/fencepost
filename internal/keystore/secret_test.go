package keystore

import (
	"testing"
)

func TestSetSecret_StoresSecret(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key123")

	if err := s.SetSecret("svc", "mysecret"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := s.GetSecret("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "mysecret" {
		t.Errorf("expected mysecret, got %q", got)
	}
}

func TestSetSecret_MissingService(t *testing.T) {
	s := newTestStore(t)
	if err := s.SetSecret("ghost", "val"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestGetSecret_MissingService(t *testing.T) {
	s := newTestStore(t)
	_, err := s.GetSecret("ghost")
	if err == nil {
		t.Error("expected error for missing service")
	}
}

func TestClearSecret_RemovesSecret(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key123")
	_ = s.SetSecret("svc", "mysecret")

	if err := s.ClearSecret("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ := s.GetSecret("svc")
	if got != "" {
		t.Errorf("expected empty secret, got %q", got)
	}
}

func TestServicesWithSecret_ReturnsMatchingSorted(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("beta", "k1")
	_ = s.Set("alpha", "k2")
	_ = s.Set("gamma", "k3")
	_ = s.SetSecret("beta", "s1")
	_ = s.SetSecret("gamma", "s2")

	result := s.ServicesWithSecret()
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
	if result[0] != "beta" || result[1] != "gamma" {
		t.Errorf("unexpected order: %v", result)
	}
}

func TestSetSecret_PersistsAcrossReload(t *testing.T) {
	path := tempStorePath(t)
	s, _ := New(path)
	_ = s.Set("svc", "key")
	_ = s.SetSecret("svc", "persisted")

	s2, _ := New(path)
	got, err := s2.GetSecret("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "persisted" {
		t.Errorf("expected persisted, got %q", got)
	}
}
