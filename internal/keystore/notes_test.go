package keystore

import (
	"testing"
)

func TestSetNote_StoresNote(t *testing.T) {
	s, err := New(tempStorePath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := s.Set("svc", "key123"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := s.SetNote("svc", "production key"); err != nil {
		t.Fatalf("SetNote: %v", err)
	}
	note, err := s.GetNote("svc")
	if err != nil {
		t.Fatalf("GetNote: %v", err)
	}
	if note != "production key" {
		t.Errorf("expected %q, got %q", "production key", note)
	}
}

func TestSetNote_MissingService(t *testing.T) {
	s, err := New(tempStorePath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := s.SetNote("ghost", "some note"); err == nil {
		t.Error("expected error for missing service, got nil")
	}
}

func TestGetNote_MissingService(t *testing.T) {
	s, err := New(tempStorePath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if _, err := s.GetNote("ghost"); err == nil {
		t.Error("expected error for missing service, got nil")
	}
}

func TestClearNote_RemovesNote(t *testing.T) {
	s, err := New(tempStorePath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := s.Set("svc", "key123"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := s.SetNote("svc", "temporary note"); err != nil {
		t.Fatalf("SetNote: %v", err)
	}
	if err := s.ClearNote("svc"); err != nil {
		t.Fatalf("ClearNote: %v", err)
	}
	note, err := s.GetNote("svc")
	if err != nil {
		t.Fatalf("GetNote after clear: %v", err)
	}
	if note != "" {
		t.Errorf("expected empty note after clear, got %q", note)
	}
}

func TestSetNote_PersistsAcrossReload(t *testing.T) {
	path := tempStorePath(t)
	s, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := s.Set123"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := s.SetNote("svc", "persisted note"); err != nil {
		t.Fatalf("SetNote: %v", err)
	}
	s2, err := New(path)
	if err != nil {
		t.Fatalf("New reload: %v", err)
	}
	note, err := s2.GetNote("svc")
	if err != nil {
		t.Fatalf("GetNote after reload: %v", err)
	}
	if note != "persisted note" {
		t.Errorf("expected %q after reload, got %q", "persisted note", note)
	}
}
