package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fencepost/internal/keystore"
)

func TestImport_StoresKey(t *testing.T) {
	dir := t.TempDir()
	storePath := filepath.Join(dir, "keys.json")

	ks, err := keystore.New(storePath)
	if err != nil {
		t.Fatalf("creating keystore: %v", err)
	}

	if err := ks.Set("myservice", "imported-key-abc123"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	got, err := ks.Get("myservice")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got != "imported-key-abc123" {
		t.Errorf("expected %q, got %q", "imported-key-abc123", got)
	}
}

func TestImport_OverwriteExistingKey(t *testing.T) {
	dir := t.TempDir()
	storePath := filepath.Join(dir, "keys.json")

	ks, err := keystore.New(storePath)
	if err != nil {
		t.Fatalf("creating keystore: %v", err)
	}

	if err := ks.Set("svc", "old-key"); err != nil {
		t.Fatalf("initial Set failed: %v", err)
	}

	if err := ks.Set("svc", "new-imported-key"); err != nil {
		t.Fatalf("overwrite Set failed: %v", err)
	}

	got, err := ks.Get("svc")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got != "new-imported-key" {
		t.Errorf("expected %q, got %q", "new-imported-key", got)
	}
}

func TestImport_PersistsAcrossReload(t *testing.T) {
	dir := t.TempDir()
	storePath := filepath.Join(dir, "keys.json")

	ks, err := keystore.New(storePath)
	if err != nil {
		t.Fatalf("creating keystore: %v", err)
	}
	if err := ks.Set("persist-svc", "persisted-key"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Reload from disk
	ks2, err := keystore.New(storePath)
	if err != nil {
		t.Fatalf("reloading keystore: %v", err)
	}
	got, err := ks2.Get("persist-svc")
	if err != nil {
		t.Fatalf("Get after reload failed: %v", err)
	}
	if got != "persisted-key" {
		t.Errorf("expected %q, got %q", "persisted-key", got)
	}

	_ = os.Remove(storePath)
}
