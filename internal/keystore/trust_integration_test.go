package keystore_test

import (
	"testing"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func TestTrust_PersistsToDisk(t *testing.T) {
	path := tempStorePath(t)
	s1 := keystore.New(path)
	s1.Set("api", "secret")
	if err := s1.SetTrustLevel("api", "full"); err != nil {
		t.Fatalf("SetTrustLevel failed: %v", err)
	}

	s2 := keystore.New(path)
	level, err := s2.GetTrustLevel("api")
	if err != nil {
		t.Fatalf("GetTrustLevel after reload failed: %v", err)
	}
	if level != "full" {
		t.Errorf("expected full after reload, got %s", level)
	}
}

func TestTrust_ClearPersistsToDisk(t *testing.T) {
	path := tempStorePath(t)
	s1 := keystore.New(path)
	s1.Set("api", "secret")
	s1.SetTrustLevel("api", "high")
	s1.ClearTrustLevel("api")

	s2 := keystore.New(path)
	level, err := s2.GetTrustLevel("api")
	if err != nil {
		t.Fatalf("GetTrustLevel after reload failed: %v", err)
	}
	if level != "none" {
		t.Errorf("expected none after clear and reload, got %s", level)
	}
}

func TestTrust_InvalidLevelRejected(t *testing.T) {
	path := tempStorePath(t)
	s := keystore.New(path)
	s.Set("api", "secret")

	err := s.SetTrustLevel("api", "absolute")
	if err == nil {
		t.Fatal("expected error for invalid trust level 'absolute'")
	}

	level, _ := s.GetTrustLevel("api")
	if level != "none" {
		t.Errorf("expected trust level to remain none, got %s", level)
	}
}
