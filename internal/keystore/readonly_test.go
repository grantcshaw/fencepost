package keystore_test

import (
	"testing"

	"github.com/densestvoid/fencepost/internal/keystore"
)

func TestSetReadOnly_StoresValue(t *testing.T) {
	s := newTestStore(t)
	mustSet(t, s, "svc", "key123")

	if err := s.SetReadOnly("svc", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := s.IsReadOnly("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got {
		t.Error("expected read-only to be true")
	}
}

func TestSetReadOnly_MissingService(t *testing.T) {
	s := newTestStore(t)
	if err := s.SetReadOnly("ghost", true); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestIsReadOnly_DefaultsFalse(t *testing.T) {
	s := newTestStore(t)
	mustSet(t, s, "svc", "key123")

	got, err := s.IsReadOnly("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got {
		t.Error("expected read-only to default to false")
	}
}

func TestIsReadOnly_MissingService(t *testing.T) {
	s := newTestStore(t)
	if _, err := s.IsReadOnly("ghost"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestReadOnlyKeys_ReturnsSorted(t *testing.T) {
	s := newTestStore(t)
	for _, svc := range []string{"zebra", "alpha", "mango"} {
		mustSet(t, s, svc, "k")
		if err := s.SetReadOnly(svc, true); err != nil {
			t.Fatalf("SetReadOnly(%q): %v", svc, err)
		}
	}
	mustSet(t, s, "other", "k")

	got := s.ReadOnlyKeys()
	want := []string{"alpha", "mango", "zebra"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestReadOnlyKeys_PersistsAcrossReload(t *testing.T) {
	path := tempStorePath(t)
	s, err := keystore.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := s.Set("svc", "key"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := s.SetReadOnly("svc", true); err != nil {
		t.Fatalf("SetReadOnly: %v", err)
	}

	s2, err := keystore.New(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	got, err := s2.IsReadOnly("svc")
	if err != nil {
		t.Fatalf("IsReadOnly: %v", err)
	}
	if !got {
		t.Error("expected read-only to persist across reload")
	}
}
