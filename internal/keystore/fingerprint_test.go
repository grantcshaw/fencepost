package keystore_test

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func TestSetFingerprint_StoresFingerprint(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	if err := s.Set("svc", "mykey"); err != nil {
		t.Fatal(err)
	}

	fp, err := s.SetFingerprint("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	hash := sha256.Sum256([]byte("mykey"))
	expected := hex.EncodeToString(hash[:])
	if fp != expected {
		t.Errorf("got %q, want %q", fp, expected)
	}
}

func TestSetFingerprint_MissingService(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	_, err := s.SetFingerprint("ghost")
	if err == nil {
		t.Error("expected error for missing service")
	}
}

func TestGetFingerprint_MissingService(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	_, err := s.GetFingerprint("ghost")
	if err == nil {
		t.Error("expected error for missing service")
	}
}

func TestClearFingerprint_RemovesFingerprint(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	if err := s.Set("svc", "mykey"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.SetFingerprint("svc"); err != nil {
		t.Fatal(err)
	}
	if err := s.ClearFingerprint("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fp, err := s.GetFingerprint("svc")
	if err != nil {
		t.Fatal(err)
	}
	if fp != "" {
		t.Errorf("expected empty fingerprint, got %q", fp)
	}
}

func TestServicesByFingerprint_ReturnsMatchingSorted(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	for _, name := range []string{"bravo", "alpha"} {
		if err := s.Set(name, "sharedkey"); err != nil {
			t.Fatal(err)
		}
		if _, err := s.SetFingerprint(name); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Set("charlie", "differentkey"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.SetFingerprint("charlie"); err != nil {
		t.Fatal(err)
	}

	hash := sha256.Sum256([]byte("sharedkey"))
	fp := hex.EncodeToString(hash[:])

	results := s.ServicesByFingerprint(fp)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0] != "alpha" || results[1] != "bravo" {
		t.Errorf("unexpected order: %v", results)
	}
}

func TestSetFingerprint_PersistsAcrossReload(t *testing.T) {
	path := tempStorePath(t)
	s := keystore.New(path)
	if err := s.Set("svc", "mykey"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.SetFingerprint("svc"); err != nil {
		t.Fatal(err)
	}

	s2 := keystore.New(path)
	hash := sha256.Sum256([]byte("mykey"))
	expected := hex.EncodeToString(hash[:])

	fp, err := s2.GetFingerprint("svc")
	if err != nil {
		t.Fatal(err)
	}
	if fp != expected {
		t.Errorf("got %q after reload, want %q", fp, expected)
	}
}
