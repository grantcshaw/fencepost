package keystore

import (
	"testing"
)

func TestSetCipher_StoresCipher(t *testing.T) {
	s := newTestStore(t)
	s.Set("svcA", "key1")

	if err := s.SetCipher("svcA", "chacha20"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := s.GetCipher("svcA")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "chacha20" {
		t.Errorf("expected chacha20, got %q", got)
	}
}

func TestSetCipher_MissingService(t *testing.T) {
	s := newTestStore(t)
	if err := s.SetCipher("ghost", "aes-256"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestSetCipher_InvalidValue(t *testing.T) {
	s := newTestStore(t)
	s.Set("svcA", "key1")
	if err := s.SetCipher("svcA", "rot13"); err == nil {
		t.Error("expected error for invalid cipher")
	}
}

func TestGetCipher_DefaultsToAES256(t *testing.T) {
	s := newTestStore(t)
	s.Set("svcA", "key1")
	got, err := s.GetCipher("svcA")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "aes-256" {
		t.Errorf("expected aes-256 default, got %q", got)
	}
}

func TestGetCipher_MissingService(t *testing.T) {
	s := newTestStore(t)
	if _, err := s.GetCipher("ghost"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestClearCipher_RemovesCipher(t *testing.T) {
	s := newTestStore(t)
	s.Set("svcA", "key1")
	s.SetCipher("svcA", "aes-128")
	if err := s.ClearCipher("svcA"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ := s.GetCipher("svcA")
	if got != "aes-256" {
		t.Errorf("expected default aes-256 after clear, got %q", got)
	}
}

func TestServicesByCipher_ReturnsMatchingSorted(t *testing.T) {
	s := newTestStore(t)
	s.Set("svcB", "k2")
	s.Set("svcA", "k1")
	s.Set("svcC", "k3")
	s.SetCipher("svcA", "chacha20")
	s.SetCipher("svcC", "chacha20")

	result, err := s.ServicesByCipher("chacha20")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 || result[0] != "svcA" || result[1] != "svcC" {
		t.Errorf("unexpected result: %v", result)
	}
}
