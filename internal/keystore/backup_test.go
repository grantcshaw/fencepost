package keystore

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)
func TestWriteBackup_CreatesFile(t *testing.T) {
	store := newTestStore(t)
	_ = store.Set("svc-a", "key-aaa")

	dir := t.TempDir()
	meta, err := store.WriteBackup(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(meta.Path); err != nil {
		t.Errorf("backup file not found: %v", err)
	}
}

func TestWriteBackup_FileNameContainsTimestamp(t *testing.T) {
	store := newTestStore(t)
	dir := t.TempDir()

	meta, err := store.WriteBackup(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	base := filepath.Base(meta.Path)
	if !strings.HasPrefix(base, "fencepost-backup-") {
		t.Errorf("unexpected file name: %s", base)
	}
}

func TestWriteBackup_ServiceCount(t *testing.T) {
	store := newTestStore(t)
	_ = store.Set("alpha", "k1")
	_ = store.Set("beta", "k2")

	dir := t.TempDir()
	meta, err := store.WriteBackup(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if meta.Services != 2 {
		t.Errorf("expected 2 services, got %d", meta.Services)
	}
}

func TestRestoreBackup_RestoresEntries(t *testing.T) {
	original := newTestStore(t)
	_ = original.Set("svc-x", "key-xyz")
	_ = original.Set("svc-y", "key-yyy")

	dir := t.TempDir()
	meta, err := original.WriteBackup(dir)
	if err != nil {
		t.Fatalf("backup failed: %v", err)
	}

	restored := newTestStore(t)
	if err := restored.RestoreBackup(meta.Path); err != nil {
		t.Fatalf("restore failed: %v", err)
	}

	for _, svc := range []string{"svc-x", "svc-y"} {
		if _, err := restored.Get(svc); err != nil {
			t.Errorf("missing service %q after restore", svc)
		}
	}
}

func TestRestoreBackup_MissingFile(t *testing.T) {
	store := newTestStore(t)
	err := store.RestoreBackup("/nonexistent/path/backup.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestWriteBackup_CreatedAtIsRecent(t *testing.T) {
	store := newTestStore(t)
	before := time.Now().UTC().Add(-time.Second)

	dir := t.TempDir()
	meta, err := store.WriteBackup(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if meta.CreatedAt.Before(before) {
		t.Errorf("Crev is before test start %v", meta.CreatedAt, before)
	}
}

func newTestStore(t *testing.T) *Store {
	t.Helper()
	path := filepath.Join(t.TempDir(), "store.json")
	s, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}
