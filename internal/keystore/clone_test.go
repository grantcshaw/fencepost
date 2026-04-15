package keystore

import (
	"testing"
	"time"
)

func TestClone_CopiesEntry(t *testing.T) {
	s := newTestStore(t)

	_ = s.Set("alpha", "key-abc")
	_ = s.SetNote("alpha", "my note")
	_, _ = s.SetTags("alpha", []string{"prod", "critical"})

	result, err := s.Clone("alpha", "alpha-copy", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Source != "alpha" || result.Destination != "alpha-copy" {
		t.Errorf("unexpected result: %+v", result)
	}
	if result.Overwritten {
		t.Error("expected Overwritten=false for new destination")
	}

	entry, err := s.Get("alpha-copy")
	if err != nil {
		t.Fatalf("cloned entry not found: %v", err)
	}
	if entry.Key != "key-abc" {
		t.Errorf("expected key %q, got %q", "key-abc", entry.Key)
	}
	if entry.Note != "my note" {
		t.Errorf("expected note %q, got %q", "my note", entry.Note)
	}
	if len(entry.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(entry.Tags))
	}
}

func TestClone_MissingSource(t *testing.T) {
	s := newTestStore(t)

	_, err := s.Clone("nonexistent", "dest", false)
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestClone_DestinationExistsNoOverwrite(t *testing.T) {
	s := newTestStore(t)

	_ = s.Set("svc1", "key-1")
	_ = s.Set("svc2", "key-2")

	_, err := s.Clone("svc1", "svc2", false)
	if err == nil {
		t.Fatal("expected error when destination exists without overwrite")
	}
}

func TestClone_OverwriteExistingDestination(t *testing.T) {
	s := newTestStore(t)

	_ = s.Set("svc1", "key-new")
	_ = s.Set("svc2", "key-old")

	result, err := s.Clone("svc1", "svc2", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Overwritten {
		t.Error("expected Overwritten=true")
	}

	entry, _ := s.Get("svc2")
	if entry.Key != "key-new" {
		t.Errorf("expected overwritten key %q, got %q", "key-new", entry.Key)
	}
}

func TestClone_TagsAreIndependent(t *testing.T) {
	s := newTestStore(t)

	_ = s.Set("origin", "key-x")
	_, _ = s.SetTags("origin", []string{"tagA"})

	_, err := s.Clone("origin", "copy", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Mutate original tags; copy should be unaffected
	_, _ = s.SetTags("origin", []string{"tagA", "tagB"})

	copyEntry, _ := s.Get("copy")
	if len(copyEntry.Tags) != 1 {
		t.Errorf("expected 1 tag on copy after mutating origin, got %d", len(copyEntry.Tags))
	}
}

func newTestStore(t *testing.T) *Store {
	t.Helper()
	path := tempStorePath(t)
	s, err := New(path, time.Now())
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	return s
}
