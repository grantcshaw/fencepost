package keystore

import (
	"testing"
)

func TestWatch_PinsAndUnwatches(t *testing.T) {
	s, err := New(tempStorePath(t))
	if err != nil {
		t.Fatal(err)
	}
	_ = s.Set("svcA", "key1")

	if err := s.Watch("svcA"); err != nil {
		t.Fatalf("Watch: %v", err)
	}
	if !s.IsWatched("svcA") {
		t.Error("expected svcA to be watched")
	}
	if err := s.Unwatch("svcA"); err != nil {
		t.Fatalf("Unwatch: %v", err)
	}
	if s.IsWatched("svcA") {
		t.Error("expected svcA to not be watched after unwatch")
	}
}

func TestWatch_MissingService(t *testing.T) {
	s, _ := New(tempStorePath(t))
	if err := s.Watch("ghost"); err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestIsWatched_MissingService(t *testing.T) {
	s, _ := New(tempStorePath(t))
	if s.IsWatched("ghost") {
		t.Error("expected false for missing service")
	}
}

func TestWatchedKeys_ReturnsSorted(t *testing.T) {
	s, _ := New(tempStorePath(t))
	_ = s.Set("beta", "k1")
	_ = s.Set("alpha", "k2")
	_ = s.Set("gamma", "k3")
	_ = s.Watch("beta")
	_ = s.Watch("alpha")
	_ = s.Watch("gamma")

	entries := s.WatchedKeys()
	if len(entries) != 3 {
		t.Fatalf("expected 3, got %d", len(entries))
	}
	if entries[0].Service != "alpha" || entries[1].Service != "beta" || entries[2].Service != "gamma" {
		t.Error("expected sorted order")
	}
}

func TestWatchedKeys_PersistsAcrossReload(t *testing.T) {
	path := tempStorePath(t)
	s, _ := New(path)
	_ = s.Set("svcX", "keyX")
	_ = s.Watch("svcX")

	s2, _ := New(path)
	if !s2.IsWatched("svcX") {
		t.Error("expected watch to persist across reload")
	}
}
