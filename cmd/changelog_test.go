package cmd_test

import (
	"testing"

	"github.com/richbl/fencepost/internal/keystore"
)

func newChangelogTestStore(t *testing.T) *keystore.Store {
	t.Helper()
	path := tempStorePath(t)
	ks, err := keystore.New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	if err := ks.Set("mysvc", "apikey-abc"); err != nil {
		t.Fatalf("failed to seed store: %v", err)
	}
	return ks
}

func TestChangelogAdd_StoresEntry(t *testing.T) {
	ks := newChangelogTestStore(t)

	if err := ks.AppendChangelog("mysvc", "rotated", "scheduled"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, err := ks.GetChangelog("mysvc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Event != "rotated" {
		t.Errorf("expected 'rotated', got %q", entries[0].Event)
	}
}

func TestChangelogClear_RemovesAllEntries(t *testing.T) {
	ks := newChangelogTestStore(t)
	_ = ks.AppendChangelog("mysvc", "rotated", "")
	_ = ks.AppendChangelog("mysvc", "imported", "bulk import")

	if err := ks.ClearChangelog("mysvc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, _ := ks.GetChangelog("mysvc")
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestChangelogList_MultipleEntries(t *testing.T) {
	ks := newChangelogTestStore(t)
	_ = ks.AppendChangelog("mysvc", "created", "initial")
	_ = ks.AppendChangelog("mysvc", "rotated", "quarterly")
	_ = ks.AppendChangelog("mysvc", "audited", "compliance check")

	entries, err := ks.GetChangelog("mysvc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[2].Event != "audited" {
		t.Errorf("expected last entry to be 'audited', got %q", entries[2].Event)
	}
}
