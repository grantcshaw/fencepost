package keystore_test

import (
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func TestMerge_EmptySource(t *testing.T) {
	dst := keystore.New(filepath.Join(t.TempDir(), "dst.json"))
	src := keystore.New(filepath.Join(t.TempDir(), "src.json"))

	results, err := dst.Merge(src, keystore.MergeOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results for empty source, got %d", len(results))
	}
}

func TestMerge_PersistsToDisk(t *testing.T) {
	dstPath := filepath.Join(t.TempDir(), "dst.json")
	srcPath := filepath.Join(t.TempDir(), "src.json")

	dst := keystore.New(dstPath)
	src := keystore.New(srcPath)
	_ = src.Set("persisted", "value-x")

	_, err := dst.Merge(src, keystore.MergeOptions{})
	if err != nil {
		t.Fatalf("merge error: %v", err)
	}

	// Reload from disk.
	reloaded := keystore.New(dstPath)
	e, err := reloaded.Get("persisted")
	if err != nil {
		t.Fatalf("entry not persisted: %v", err)
	}
	if e.Key != "value-x" {
		t.Errorf("expected value-x, got %s", e.Key)
	}
}

func TestMerge_MixedActions(t *testing.T) {
	dstPath := filepath.Join(t.TempDir(), "dst.json")
	srcPath := filepath.Join(t.TempDir(), "src.json")

	dst := keystore.New(dstPath)
	src := keystore.New(srcPath)
	_ = dst.Set("existing", "old")
	_ = src.Set("existing", "new")
	_ = src.Set("fresh", "brand-new")

	results, err := dst.Merge(src, keystore.MergeOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("merge error: %v", err)
	}

	actionFor := map[string]string{}
	for _, r := range results {
		actionFor[r.Service] = r.Action
	}
	if actionFor["existing"] != "overwritten" {
		t.Errorf("expected overwritten, got %q", actionFor["existing"])
	}
	if actionFor["fresh"] != "added" {
		t.Errorf("expected added, got %q", actionFor["fresh"])
	}
}
