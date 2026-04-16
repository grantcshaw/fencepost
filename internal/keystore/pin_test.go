package keystore

import (
	"testing"
)

func TestPin_PinsAndUnpins(t *testing.T) {
	s, _ := New(tempStorePath(t))
	_ = s.Set("svcA", "key1")

	if err := s.PinEntry("svcA"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	pinned, err := s.IsPinned("svcA")
	if err != nil || !pinned {
		t.Fatalf("expected pinned=true, got %v (err: %v)", pinned, err)
	}

	_ = s.UnpinEntry("svcA")
	pinned, _ = s.IsPinned("svcA")
	if pinned {
		t.Fatal("expected pinned=false after unpin")
	}
}

func TestPin_MissingService(t *testing.T) {
	s, _ := New(tempStorePath(t))
	if err := s.PinEntry("ghost"); err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestIsPinned_MissingService(t *testing.T) {
	s, _ := New(tempStorePath(t))
	_, err := s.IsPinned("ghost")
	if err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestPinnedKeys_ReturnsSorted(t *testing.T) {
	s, _ := New(tempStorePath(t))
	_ = s.Set("zebra", "k1")
	_ = s.Set("alpha", "k2")
	_ = s.Set("mango", "k3")
	_ = s.PinEntry("zebra")
	_ = s.PinEntry("alpha")

	pinned := s.PinnedKeys()
	if len(pinned) != 2 {
		t.Fatalf("expected 2 pinned, got %d", len(pinned))
	}
	if pinned[0] != "alpha" || pinned[1] != "zebra" {
		t.Fatalf("unexpected order: %v", pinned)
	}
}

func TestPinnedKeys_PersistsAcrossReload(t *testing.T) {
	path := tempStorePath(t)
	s, _ := New(path)
	_ = s.Set("svc1", "key")
	_ = s.PinEntry("svc1")

	s2, _ := New(path)
	pinned, err := s2.IsPinned("svc1")
	if err != nil || !pinned {
		t.Fatalf("expected pinned to persist, got %v (err: %v)", pinned, err)
	}
}
