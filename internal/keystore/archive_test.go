package keystore

import (
	"os"
	"path/filepath"
	"testing"
)

func tempArchivePath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "archive.json")
}

func TestArchive_RemovesFromStore(t *testing.T) {
	store := newTestStore(t)
	_ = store.Set("mysvc", "key123")
	archPath := tempArchivePath(t)

	if err := store.Archive("mysvc", archPath); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := store.Get("mysvc"); err == nil {
		t.Error("expected service to be removed from active store")
	}
}

func TestArchive_WritesEntryToFile(t *testing.T) {
	store := newTestStore(t)
	_ = store.Set("svc-a", "secret")
	archPath := tempArchivePath(t)

	if err := store.Archive("svc-a", archPath); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, err := LoadArchive(archPath)
	if err != nil {
		t.Fatalf("LoadArchive: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 archived entry, got %d", len(entries))
	}
	if entries[0].Entry.Key != "secret" {
		t.Errorf("expected key %q, got %q", "secret", entries[0].Entry.Key)
	}
	if entries[0].ArchivedAt.IsZero() {
		t.Error("expected ArchivedAt to be set")
	}
}

func TestArchive_AppendsToExistingArchive(t *testing.T) {
	store := newTestStore(t)
	_ = store.Set("svc-1", "key1")
	_ = store.Set("svc-2", "key2")
	archPath := tempArchivePath(t)

	_ = store.Archive("svc-1", archPath)
	_ = store.Archive("svc-2", archPath)

	entries, err := LoadArchive(archPath)
	if err != nil {
		t.Fatalf("LoadArchive: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 archived entries, got %d", len(entries))
	}
}

func TestArchive_MissingService(t *testing.T) {
	store := newTestStore(t)
	archPath := tempArchivePath(t)

	err := store.Archive("nonexistent", archPath)
	if err == nil {
		t.Error("expected error for missing service")
	}
}

func TestLoadArchive_MissingFile(t *testing.T) {
	entries, err := LoadArchive(filepath.Join(t.TempDir(), "no-such-file.json"))
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(entries))
	}
}

// newTestStore is a helper shared across keystore tests.
func newTestStore(t *testing.T) *Store {
	t.Helper()
	path := filepath.Join(t.TempDir(), "store.json")
	store, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_ = os.WriteFile // keep import if needed
	return store
}
