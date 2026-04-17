package keystore

import (
	"testing"
)

func TestSetVersion_StoresVersion(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key123")

	if err := s.SetVersion("svc", "v2"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, err := s.GetVersion("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "v2" {
		t.Errorf("expected v2, got %q", v)
	}
}

func TestSetVersion_MissingService(t *testing.T) {
	s := newTestStore(t)
	if err := s.SetVersion("ghost", "v1"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestGetVersion_MissingService(t *testing.T) {
	s := newTestStore(t)
	if _, err := s.GetVersion("ghost"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestClearVersion_RemovesVersion(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key123")
	_ = s.SetVersion("svc", "v3")

	if err := s.ClearVersion("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, _ := s.GetVersion("svc")
	if v != "" {
		t.Errorf("expected empty version, got %q", v)
	}
}

func TestServicesByVersion_ReturnsMatchingSorted(t *testing.T) {
	s := newTestStore(t)
	for _, svc := range []string{"gamma", "alpha", "beta"} {
		_ = s.Set(svc, "k")
		_ = s.SetVersion(svc, "v1")
	}
	_ = s.Set("delta", "k")
	_ = s.SetVersion("delta", "v2")

	result := s.ServicesByVersion("v1")
	if len(result) != 3 {
		t.Fatalf("expected 3, got %d", len(result))
	}
	if result[0] != "alpha" || result[1] != "beta" || result[2] != "gamma" {
		t.Errorf("unexpected order: %v", result)
	}
}

func TestSetVersion_PersistsAcrossReload(t *testing.T) {
	path := tempStorePath(t)
	s, _ := New(path)
	_ = s.Set("svc", "key")
	_ = s.SetVersion("svc", "v5")

	s2, _ := New(path)
	v, err := s2.GetVersion("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "v5" {
		t.Errorf("expected v5, got %q", v)
	}
}
