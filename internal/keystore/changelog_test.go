package keystore_test

import (
	"testing"

	"github.com/richbl/fencepost/internal/keystore"
)

func TestAppendChangelog_StoresEntry(t *testing.T) {
	store := newTestStore(t)
	_ = store.Set("svc", "key123")

	if err := store.AppendChangelog("svc", "rotated", "manual rotation"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, err := store.GetChangelog("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Event != "rotated" {
		t.Errorf("expected event 'rotated', got %q", entries[0].Event)
	}
	if entries[0].Detail != "manual rotation" {
		t.Errorf("expected detail 'manual rotation', got %q", entries[0].Detail)
	}
}

func TestAppendChangelog_MissingService(t *testing.T) {
	store := newTestStore(t)
	err := store.AppendChangelog("ghost", "rotated", "")
	if err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestGetChangelog_MissingService(t *testing.T) {
	store := newTestStore(t)
	_, err := store.GetChangelog("ghost")
	if err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestClearChangelog_RemovesEntries(t *testing.T) {
	store := newTestStore(t)
	_ = store.Set("svc", "key123")
	_ = store.AppendChangelog("svc", "rotated", "")
	_ = store.AppendChangelog("svc", "imported", "")

	if err := store.ClearChangelog("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, _ := store.GetChangelog("svc")
	if len(entries) != 0 {
		t.Errorf("expected 0 entries after clear, got %d", len(entries))
	}
}

func TestServicesWithChangelog_ReturnsMatchingSorted(t *testing.T) {
	store := newTestStore(t)
	_ = store.Set("beta", "k1")
	_ = store.Set("alpha", "k2")
	_ = store.Set("gamma", "k3")

	_ = store.AppendChangelog("beta", "rotated", "")
	_ = store.AppendChangelog("gamma", "imported", "")

	result := store.ServicesWithChangelog()
	if len(result) != 2 {
		t.Fatalf("expected 2 services, got %d", len(result))
	}
	if result[0] != "beta" || result[1] != "gamma" {
		t.Errorf("unexpected order: %v", result)
	}
}

func TestAppendChangelog_PersistsAcrossReload(t *testing.T) {
	path := tempStorePath(t)
	store, _ := keystore.New(path)
	_ = store.Set("svc", "key")
	_ = store.AppendChangelog("svc", "created", "initial import")

	reloaded, _ := keystore.New(path)
	entries, err := reloaded.GetChangelog("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 || entries[0].Event != "created" {
		t.Errorf("changelog did not persist: %+v", entries)
	}
}
