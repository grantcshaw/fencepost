package keystore

import (
	"testing"
)

func TestSetComment_StoresComment(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key123")

	if err := s.SetComment("svc", "primary prod key"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := s.GetComment("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "primary prod key" {
		t.Errorf("expected %q, got %q", "primary prod key", got)
	}
}

func TestSetComment_MissingService(t *testing.T) {
	s := newTestStore(t)
	if err := s.SetComment("ghost", "nope"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestGetComment_MissingService(t *testing.T) {
	s := newTestStore(t)
	if _, err := s.GetComment("ghost"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestClearComment_RemovesComment(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key123")
	_ = s.SetComment("svc", "temp comment")

	if err := s.ClearComment("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, _ := s.GetComment("svc")
	if got != "" {
		t.Errorf("expected empty comment, got %q", got)
	}
}

func TestSetComment_PersistsAcrossReload(t *testing.T) {
	path := tempStorePath(t)
	s, _ := New(path)
	_ = s.Set("svc", "key123")
	_ = s.SetComment("svc", "persisted comment")

	s2, _ := New(path)
	got, err := s2.GetComment("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "persisted comment" {
		t.Errorf("expected %q, got %q", "persisted comment", got)
	}
}
