package keystore

import (
	"testing"
)

func TestCopyKey_CopiesEntry(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("alpha", "key-alpha")

	if err := s.CopyKey("alpha", "beta", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	v, err := s.Get("beta")
	if err != nil || v != "key-alpha" {
		t.Fatalf("expected key-alpha, got %q (err %v)", v, err)
	}

	// source still present
	v2, err := s.Get("alpha")
	if err != nil || v2 != "key-alpha" {
		t.Fatalf("source should still exist, got %q (err %v)", v2, err)
	}
}

func TestCopyKey_MissingSource(t *testing.T) {
	s := newTestStore(t)
	if err := s.CopyKey("missing", "dst", false); err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestCopyKey_DestinationExistsNoOverwrite(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("alpha", "key-alpha")
	_ = s.Set("beta", "key-beta")

	if err := s.CopyKey("alpha", "beta", false); err == nil {
		t.Fatal("expected error when destination exists without overwrite")
	}
}

func TestCopyKey_OverwriteExistingDestination(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("alpha", "key-alpha")
	_ = s.Set("beta", "key-beta")

	if err := s.CopyKey("alpha", "beta", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	v, _ := s.Get("beta")
	if v != "key-alpha" {
		t.Fatalf("expected key-alpha after overwrite, got %q", v)
	}
}

func TestCopyKey_TagsAreIndependent(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("alpha", "key-alpha")
	_, _ = s.SetTags("alpha", []string{"prod"})

	_ = s.CopyKey("alpha", "beta", false)
	_, _ = s.SetTags("beta", []string{"staging"})

	tags, _ := s.GetTags("alpha")
	if len(tags) != 1 || tags[0] != "prod" {
		t.Fatalf("source tags mutated: %v", tags)
	}
}
