package keystore

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWriteSnapshot_CreatesFile(t *testing.T) {
	store := New(tempStorePath(t))
	_ = store.Set("alpha", "key-abc", time.Now())
	_ = store.Set("beta", "key-xyz", time.Now())

	dir := t.TempDir()
	path, err := store.WriteSnapshot(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("snapshot file not found: %v", err)
	}
}

func TestWriteSnapshot_ContainsAllEntries(t *testing.T) {
	store := New(tempStorePath(t))
	_ = store.Set("svc-a", "key-1", time.Now())
	_ = store.Set("svc-b", "key-2", time.Now())

	dir := t.TempDir()
	path, err := store.WriteSnapshot(dir)
	if err != nil {
		t.Fatalf("WriteSnapshot error: %v", err)
	}

	snap, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot error: %v", err)
	}

	if len(snap.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(snap.Entries))
	}
	if snap.Entries["svc-a"].Key != "key-1" {
		t.Errorf("expected key-1 for svc-a, got %q", snap.Entries["svc-a"].Key)
	}
}

func TestWriteSnapshot_EmptyStore(t *testing.T) {
	store := New(tempStorePath(t))
	dir := t.TempDir()

	path, err := store.WriteSnapshot(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	snap, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot error: %v", err)
	}
	if len(snap.Entries) != 0 {
		t.Errorf("expected empty snapshot, got %d entries", len(snap.Entries))
	}
}

func TestLoadSnapshot_MissingFile(t *testing.T) {
	_, err := LoadSnapshot(filepath.Join(t.TempDir(), "nonexistent.json"))
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestWriteSnapshot_FileNameContainsTimestamp(t *testing.T) {
	store := New(tempStorePath(t))
	dir := t.TempDir()

	before := time.Now().Unix()
	path, err := store.WriteSnapshot(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	after := time.Now().Unix()

	base := filepath.Base(path)
	var ts int64
	if _, err := parseSnapshotTimestamp(base, &ts); err != nil {
		t.Fatalf("could not parse timestamp from filename %q: %v", base, err)
	}
	if ts < before || ts > after {
		t.Errorf("timestamp %d out of range [%d, %d]", ts, before, after)
	}
}

func parseSnapshotTimestamp(name string, out *int64) (int, error) {
	_, err := fmt.Sscanf(name, "snapshot-%d.json", out)
	return 0, err
}
