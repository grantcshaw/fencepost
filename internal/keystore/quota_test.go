package keystore

import (
	"testing"
)

func TestSetQuota_StoresQuota(t *testing.T) {
	s := newTestStore(t)
	if err := s.Set("svc", "key"); err != nil {
		t.Fatal(err)
	}
	if err := s.SetQuota("svc", 1000); err != nil {
		t.Fatal(err)
	}
	q, err := s.GetQuota("svc")
	if err != nil {
		t.Fatal(err)
	}
	if q != 1000 {
		t.Fatalf("expected 1000, got %d", q)
	}
}

func TestSetQuota_MissingService(t *testing.T) {
	s := newTestStore(t)
	if err := s.SetQuota("ghost", 500); err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestGetQuota_MissingService(t *testing.T) {
	s := newTestStore(t)
	if _, err := s.GetQuota("ghost"); err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestClearQuota_RemovesQuota(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key")
	_ = s.SetQuota("svc", 999)
	if err := s.ClearQuota("svc"); err != nil {
		t.Fatal(err)
	}
	q, _ := s.GetQuota("svc")
	if q != 0 {
		t.Fatalf("expected 0 after clear, got %d", q)
	}
}

func TestServicesByQuota_ReturnsMatchingSorted(t *testing.T) {
	s := newTestStore(t)
	for _, svc := range []string{"alpha", "beta", "gamma"} {
		_ = s.Set(svc, "k")
	}
	_ = s.SetQuota("alpha", 100)
	_ = s.SetQuota("beta", 500)
	_ = s.SetQuota("gamma", 50)

	result := s.ServicesByQuota(100)
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
	if result[0] != "alpha" || result[1] != "beta" {
		t.Fatalf("unexpected order: %v", result)
	}
}

func TestSetQuota_PersistsAcrossReload(t *testing.T) {
	path := tempStorePath(t)
	s, _ := New(path)
	_ = s.Set("svc", "key")
	_ = s.SetQuota("svc", 250)

	s2, _ := New(path)
	q, err := s2.GetQuota("svc")
	if err != nil {
		t.Fatal(err)
	}
	if q != 250 {
		t.Fatalf("expected 250 after reload, got %d", q)
	}
}
