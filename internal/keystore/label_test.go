package keystore

import (
	"testing"
	"time"
)

func TestSetLabel_StoresLabel(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key123", time.Now())

	if err := s.SetLabel("svc", "My Service"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	label, err := s.GetLabel("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if label != "My Service" {
		t.Errorf("expected 'My Service', got %q", label)
	}
}

func TestSetLabel_MissingService(t *testing.T) {
	s := newTestStore(t)
	if err := s.SetLabel("ghost", "label"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestGetLabel_MissingService(t *testing.T) {
	s := newTestStore(t)
	_, err := s.GetLabel("ghost")
	if err == nil {
		t.Error("expected error for missing service")
	}
}

func TestClearLabel_RemovesLabel(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key123", time.Now())
	_ = s.SetLabel("svc", "Temp Label")

	if err := s.ClearLabel("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	label, _ := s.GetLabel("svc")
	if label != "" {
		t.Errorf("expected empty label, got %q", label)
	}
}

func TestLabeledKeys_ReturnsSorted(t *testing.T) {
	s := newTestStore(t)
	for _, svc := range []string{"zebra", "alpha", "middle"} {
		_ = s.Set(svc, "k", time.Now())
	}
	_ = s.SetLabel("zebra", "Z")
	_ = s.SetLabel("alpha", "A")

	keys := s.LabeledKeys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 labeled keys, got %d", len(keys))
	}
	if keys[0] != "alpha" || keys[1] != "zebra" {
		t.Errorf("unexpected order: %v", keys)
	}
}

func TestSetLabel_PersistsAcrossReload(t *testing.T) {
	path := tempStorePath(t)
	s, _ := New(path)
	_ = s.Set("svc", "key", time.Now())
	_ = s.SetLabel("svc", "Persistent")

	s2, _ := New(path)
	label, err := s2.GetLabel("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if label != "Persistent" {
		t.Errorf("expected 'Persistent', got %q", label)
	}
}
