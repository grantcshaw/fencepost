package keystore

import (
	"testing"
	"time"
)

func TestList_EmptyStore(t *testing.T) {
	path := tempStorePath(t)
	ks, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	got := ks.List()
	if len(got) != 0 {
		t.Errorf("expected empty list, got %v", got)
	}
}

func TestList_ReturnsSortedNames(t *testing.T) {
	path := tempStorePath(t)
	ks, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	services := []string{"zebra", "alpha", "mango", "beta"}
	for _, svc := range services {
		if err := ks.Set(svc, "key-"+svc); err != nil {
			t.Fatalf("Set(%s): %v", svc, err)
		}
	}

	got := ks.List()
	want := []string{"alpha", "beta", "mango", "zebra"}
	if len(got) != len(want) {
		t.Fatalf("expected %d entries, got %d", len(want), len(got))
	}
	for i, name := range want {
		if got[i] != name {
			t.Errorf("index %d: expected %q, got %q", i, name, got[i])
		}
	}
}

func TestList_ReflectsAfterRotate(t *testing.T) {
	path := tempStorePath(t)
	ks, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if err := ks.Set("svc-a", "original-key"); err != nil {
		t.Fatalf("Set: %v", err)
	}

	if err := ks.Rotate("svc-a", "new-key", time.Now()); err != nil {
		t.Fatalf("Rotate: %v", err)
	}

	got := ks.List()
	if len(got) != 1 || got[0] != "svc-a" {
		t.Errorf("expected [svc-a], got %v", got)
	}
}
