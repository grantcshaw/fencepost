package cmd

import (
	"testing"

	"github.com/fencepost/internal/keystore"
)

func newCipherTestStore(t *testing.T) *keystore.Store {
	t.Helper()
	path := t.TempDir() + "/store.json"
	s, err := keystore.New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	return s
}

func TestCipherCmd_SetAndGet(t *testing.T) {
	s := newCipherTestStore(t)
	if err := s.Set("myservice", "supersecret"); err != nil {
		t.Fatalf("Set: %v", err)
	}

	if err := s.SetCipher("myservice", "chacha20"); err != nil {
		t.Fatalf("SetCipher: %v", err)
	}
	got, err := s.GetCipher("myservice")
	if err != nil {
		t.Fatalf("GetCipher: %v", err)
	}
	if got != "chacha20" {
		t.Errorf("expected chacha20, got %q", got)
	}
}

func TestCipherCmd_DefaultCipher(t *testing.T) {
	s := newCipherTestStore(t)
	s.Set("svc", "key")

	got, err := s.GetCipher("svc")
	if err != nil {
		t.Fatalf("GetCipher: %v", err)
	}
	if got != "aes-256" {
		t.Errorf("expected default aes-256, got %q", got)
	}
}

func TestCipherCmd_ListByCipher(t *testing.T) {
	s := newCipherTestStore(t)
	s.Set("alpha", "k1")
	s.Set("beta", "k2")
	s.Set("gamma", "k3")
	s.SetCipher("alpha", "none")
	s.SetCipher("gamma", "none")

	result, err := s.ServicesByCipher("none")
	if err != nil {
		t.Fatalf("ServicesByCipher: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 results, got %d", len(result))
	}
	if result[0] != "alpha" || result[1] != "gamma" {
		t.Errorf("unexpected order: %v", result)
	}
}
