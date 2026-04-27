package keystore_test

import (
	"testing"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func TestPurpose_PersistsToDisk(t *testing.T) {
	path := tempStorePath(t)

	s1, err := keystore.New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	s1.Set("svc", "apikey")
	if err := s1.SetPurpose("svc", "integration"); err != nil {
		t.Fatalf("SetPurpose failed: %v", err)
	}

	s2, err := keystore.New(path)
	if err != nil {
		t.Fatalf("failed to reload store: %v", err)
	}
	got, err := s2.GetPurpose("svc")
	if err != nil {
		t.Fatalf("GetPurpose failed: %v", err)
	}
	if got != "integration" {
		t.Errorf("expected integration after reload, got %q", got)
	}
}

func TestPurpose_ClearPersistsToDisk(t *testing.T) {
	path := tempStorePath(t)

	s1, err := keystore.New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	s1.Set("svc", "apikey")
	s1.SetPurpose("svc", "webhook")
	if err := s1.ClearPurpose("svc"); err != nil {
		t.Fatalf("ClearPurpose failed: %v", err)
	}

	s2, err := keystore.New(path)
	if err != nil {
		t.Fatalf("failed to reload store: %v", err)
	}
	got, err := s2.GetPurpose("svc")
	if err != nil {
		t.Fatalf("GetPurpose failed: %v", err)
	}
	if got != "unknown" {
		t.Errorf("expected unknown after clear+reload, got %q", got)
	}
}

func TestPurpose_InvalidValueRejected(t *testing.T) {
	path := tempStorePath(t)
	s, err := keystore.New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	s.Set("svc", "apikey")
	if err := s.SetPurpose("svc", "notreal"); err == nil {
		t.Error("expected error for invalid purpose value")
	}
}
