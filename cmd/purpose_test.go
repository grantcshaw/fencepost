package cmd_test

import (
	"testing"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func newPurposeTestStore(t *testing.T) *keystore.Store {
	t.Helper()
	p := tempStorePath(t)
	s, err := keystore.New(p)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	return s
}

func TestPurposeCmd_SetAndGet(t *testing.T) {
	s := newPurposeTestStore(t)
	s.Set("myapi", "secretkey")

	if err := s.SetPurpose("myapi", "authentication"); err != nil {
		t.Fatalf("SetPurpose failed: %v", err)
	}

	got, err := s.GetPurpose("myapi")
	if err != nil {
		t.Fatalf("GetPurpose failed: %v", err)
	}
	if got != "authentication" {
		t.Errorf("expected authentication, got %q", got)
	}
}

func TestPurposeCmd_DefaultPurpose(t *testing.T) {
	s := newPurposeTestStore(t)
	s.Set("myapi", "secretkey")

	got, err := s.GetPurpose("myapi")
	if err != nil {
		t.Fatalf("GetPurpose failed: %v", err)
	}
	if got != "unknown" {
		t.Errorf("expected unknown default, got %q", got)
	}
}

func TestPurposeCmd_ListByPurpose(t *testing.T) {
	s := newPurposeTestStore(t)
	s.Set("svc-a", "k1")
	s.Set("svc-b", "k2")
	s.Set("svc-c", "k3")
	s.SetPurpose("svc-a", "signing")
	s.SetPurpose("svc-b", "signing")
	s.SetPurpose("svc-c", "encryption")

	results, err := s.ServicesByPurpose("signing")
	if err != nil {
		t.Fatalf("ServicesByPurpose failed: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 signing services, got %d", len(results))
	}
}
