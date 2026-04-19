package keystore

import (
	"testing"
)

func TestSetMetadata_StoresValue(t *testing.T) {
	s := newTestStore(t)
	mustSet(t, s, "svc", "key123")

	if err := s.SetMetadata("svc", "env", "production"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	val, err := s.GetMetadata("svc", "env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "production" {
		t.Errorf("expected production, got %q", val)
	}
}

func TestSetMetadata_MissingService(t *testing.T) {
	s := newTestStore(t)
	if err := s.SetMetadata("ghost", "k", "v"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestGetMetadata_MissingService(t *testing.T) {
	s := newTestStore(t)
	_, err := s.GetMetadata("ghost", "k")
	if err == nil {
		t.Error("expected error for missing service")
	}
}

func TestClearMetadata_RemovesKey(t *testing.T) {
	s := newTestStore(t)
	mustSet(t, s, "svc", "key123")
	_ = s.SetMetadata("svc", "env", "staging")
	if err := s.ClearMetadata("svc", "env"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	val, _ := s.GetMetadata("svc", "env")
	if val != "" {
		t.Errorf("expected empty after clear, got %q", val)
	}
}

func TestAllMetadata_ReturnsAllFields(t *testing.T) {
	s := newTestStore(t)
	mustSet(t, s, "svc", "key123")
	_ = s.SetMetadata("svc", "env", "prod")
	_ = s.SetMetadata("svc", "team", "platform")

	m, err := s.AllMetadata("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["env"] != "prod" || m["team"] != "platform" {
		t.Errorf("unexpected metadata map: %v", m)
	}
}

func TestSetMetadata_PersistsAcrossReload(t *testing.T) {
	path := tempStorePath(t)
	s, _ := New(path)
	mustSet(t, s, "svc", "key123")
	_ = s.SetMetadata("svc", "region", "us-east-1")

	s2, _ := New(path)
	val, err := s2.GetMetadata("svc", "region")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "us-east-1" {
		t.Errorf("expected us-east-1, got %q", val)
	}
}
