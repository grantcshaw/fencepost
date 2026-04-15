package keystore

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestExportAll_WritesAllServices(t *testing.T) {
	store := New(tempStorePath(t))
	_ = store.Set("alpha", "key-a")
	_ = store.Set("beta", "key-b")

	out := filepath.Join(t.TempDir(), "export.json")
	if err := store.ExportAll(out); err != nil {
		t.Fatalf("ExportAll: %v", err)
	}

	data, _ := os.ReadFile(out)
	var entries []ExportEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	// sorted order
	if entries[0].Service != "alpha" || entries[1].Service != "beta" {
		t.Errorf("unexpected order: %v", entries)
	}
}

func TestExportAll_IncludesTagsAndNote(t *testing.T) {
	store := New(tempStorePath(t))
	_ = store.Set("svc", "my-key")
	_, _ = store.SetTags("svc", []string{"prod", "critical"})
	_, _ = store.SetNote("svc", "important service")

	out := filepath.Join(t.TempDir(), "export.json")
	_ = store.ExportAll(out)

	data, _ := os.ReadFile(out)
	var entries []ExportEntry
	_ = json.Unmarshal(data, &entries)

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry")
	}
	if entries[0].Note != "important service" {
		t.Errorf("note not exported: %q", entries[0].Note)
	}
	if len(entries[0].Tags) != 2 {
		t.Errorf("tags not exported: %v", entries[0].Tags)
	}
}

func TestExportServices_SubsetOnly(t *testing.T) {
	store := New(tempStorePath(t))
	_ = store.Set("alpha", "key-a")
	_ = store.Set("beta", "key-b")
	_ = store.Set("gamma", "key-c")

	out := filepath.Join(t.TempDir(), "export.json")
	if err := store.ExportServices(out, []string{"alpha", "gamma"}); err != nil {
		t.Fatalf("ExportServices: %v", err)
	}

	data, _ := os.ReadFile(out)
	var entries []ExportEntry
	_ = json.Unmarshal(data, &entries)

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Service != "alpha" || entries[1].Service != "gamma" {
		t.Errorf("unexpected entries: %v", entries)
	}
}

func TestExportServices_MissingServiceReturnsError(t *testing.T) {
	store := New(tempStorePath(t))
	_ = store.Set("alpha", "key-a")

	out := filepath.Join(t.TempDir(), "export.json")
	err := store.ExportServices(out, []string{"alpha", "missing"})
	if err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestExportAll_EmptyStore(t *testing.T) {
	store := New(tempStorePath(t))

	out := filepath.Join(t.TempDir(), "export.json")
	if err := store.ExportAll(out); err != nil {
		t.Fatalf("ExportAll on empty store: %v", err)
	}

	data, _ := os.ReadFile(out)
	var entries []ExportEntry
	_ = json.Unmarshal(data, &entries)
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}
