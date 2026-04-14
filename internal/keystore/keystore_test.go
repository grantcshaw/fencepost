package keystore

import (
	"os"
	"path/filepath"
	"testing"
)

func tempStorePath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "keys.json")
}

func TestNew_CreatesEmptyStore(t *testing.T) {
	path := tempStorePath(t)
	s, err := New(path)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	if len(s.List()) != 0 {
		t.Errorf("expected empty store, got %d entries", len(s.List()))
	}
}

func TestStore_SetAndGet(t *testing.T) {
	s, _ := New(tempStorePath(t))

	if err := s.Set("github", "ghp_testkey123"); err != nil {
		t.Fatalf("Set() error: %v", err)
	}

	entry, err := s.Get("github")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if entry.Key != "ghp_testkey123" {
		t.Errorf("expected key %q, got %q", "ghp_testkey123", entry.Key)
	}
	if entry.Service != "github" {
		t.Errorf("expected service %q, got %q", "github", entry.Service)
	}
}

func TestStore_Rotate_UpdatesRotatedAt(t *testing.T) {
	s, _ := New(tempStorePath(t))
	s.Set("stripe", "sk_old")

	first, _ := s.Get("stripe")
	s.Set("stripe", "sk_new")
	second, _ := s.Get("stripe")

	if second.Key != "sk_new" {
		t.Errorf("expected rotated key, got %q", second.Key)
	}
	if second.RotatedAt.IsZero() {
		t.Error("expected RotatedAt to be set after rotation")
	}
	if !second.CreatedAt.Equal(first.CreatedAt) {
		t.Error("expected CreatedAt to remain unchanged after rotation")
	}
}

func TestStore_GetMissing(t *testing.T) {
	s, _ := New(tempStorePath(t))
	_, err := s.Get("nonexistent")
	if err != ErrKeyNotFound {
		t.Errorf("expected ErrKeyNotFound, got %v", err)
	}
}

func TestStore_Delete(t *testing.T) {
	s, _ := New(tempStorePath(t))
	s.Set("aws", "AKIAIOSFODNN7EXAMPLE")

	if err := s.Delete("aws"); err != nil {
		t.Fatalf("Delete() error: %v", err)
	}
	_, err := s.Get("aws")
	if err != ErrKeyNotFound {
		t.Errorf("expected ErrKeyNotFound after delete, got %v", err)
	}
}

func TestStore_Persistence(t *testing.T) {
	path := tempStorePath(t)
	s1, _ := New(path)
	s1.Set("datadog", "dd_api_key_abc")

	s2, err := New(path)
	if err != nil {
		t.Fatalf("reload error: %v", err)
	}
	entry, err := s2.Get("datadog")
	if err != nil {
		t.Fatalf("Get() after reload error: %v", err)
	}
	if entry.Key != "dd_api_key_abc" {
		t.Errorf("expected persisted key, got %q", entry.Key)
	}
}

func TestNew_MissingFile_NoError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing", "keys.json")
	_, err := New(path)
	if err == nil {
		t.Error("expected error for unwritable path, got nil")
	}
	_ = os.MkdirAll(filepath.Dir(path), 0755)
}
