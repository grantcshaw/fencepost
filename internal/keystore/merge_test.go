package keystore

import (
	"testing"
)

func TestMerge_AddsNewServices(t *testing.T) {
	dst := New(tempStorePath(t))
	src := New(tempStorePath(t))
	_ = src.Set("alpha", "key-alpha")
	_ = src.Set("beta", "key-beta")

	results, err := dst.Merge(src, MergeOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Action != "added" {
			t.Errorf("expected added, got %q for %s", r.Action, r.Service)
		}
	}
	e, _ := dst.Get("alpha")
	if e.Key != "key-alpha" {
		t.Errorf("expected key-alpha, got %s", e.Key)
	}
}

func TestMerge_SkipsExistingWithoutOverwrite(t *testing.T) {
	dst := New(tempStorePath(t))
	src := New(tempStorePath(t))
	_ = dst.Set("alpha", "original")
	_ = src.Set("alpha", "new-key")

	results, err := dst.Merge(src, MergeOptions{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Action != "skipped" {
		t.Errorf("expected skipped, got %q", results[0].Action)
	}
	e, _ := dst.Get("alpha")
	if e.Key != "original" {
		t.Errorf("key should not have changed")
	}
}

func TestMerge_OverwritesExisting(t *testing.T) {
	dst := New(tempStorePath(t))
	src := New(tempStorePath(t))
	_ = dst.Set("alpha", "original")
	_ = src.Set("alpha", "new-key")

	results, err := dst.Merge(src, MergeOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Action != "overwritten" {
		t.Errorf("expected overwritten, got %q", results[0].Action)
	}
	e, _ := dst.Get("alpha")
	if e.Key != "new-key" {
		t.Errorf("expected new-key, got %s", e.Key)
	}
}

func TestMerge_CarriesTagsAndNote(t *testing.T) {
	dst := New(tempStorePath(t))
	src := New(tempStorePath(t))
	_ = src.Set("svc", "key1")
	_ = src.SetTags("svc", []string{"prod", "critical"})
	_ = src.SetNote("svc", "important")

	_, err := dst.Merge(src, MergeOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tags, _ := dst.GetTags("svc")
	if len(tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(tags))
	}
	note, _ := dst.GetNote("svc")
	if note != "important" {
		t.Errorf("expected note 'important', got %q", note)
	}
}
