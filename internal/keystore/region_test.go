package keystore

import (
	"testing"
)

func TestSetRegion_StoresRegion(t *testing.T) {
	s := newTestStore(t)
	s.Set("svc", "key123")
	if err := s.SetRegion("svc", "us-east-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	region, err := s.GetRegion("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if region != "us-east-1" {
		t.Errorf("expected us-east-1, got %q", region)
	}
}

func TestSetRegion_MissingService(t *testing.T) {
	s := newTestStore(t)
	if err := s.SetRegion("missing", "us-west-2"); err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestGetRegion_MissingService(t *testing.T) {
	s := newTestStore(t)
	if _, err := s.GetRegion("missing"); err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestClearRegion_RemovesRegion(t *testing.T) {
	s := newTestStore(t)
	s.Set("svc", "key123")
	s.SetRegion("svc", "eu-central-1")
	if err := s.ClearRegion("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	region, _ := s.GetRegion("svc")
	if region != "" {
		t.Errorf("expected empty region, got %q", region)
	}
}

func TestServicesByRegion_ReturnsMatchingSorted(t *testing.T) {
	s := newTestStore(t)
	for _, svc := range []string{"charlie", "alpha", "beta"} {
		s.Set(svc, "k")
		s.SetRegion(svc, "ap-southeast-1")
	}
	s.Set("other", "k")
	s.SetRegion("other", "us-east-1")

	results := s.ServicesByRegion("ap-southeast-1")
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if results[0] != "alpha" || results[1] != "beta" || results[2] != "charlie" {
		t.Errorf("unexpected order: %v", results)
	}
}

func TestSetRegion_PersistsAcrossReload(t *testing.T) {
	path := tempStorePath(t)
	s, _ := New(path)
	s.Set("svc", "key123")
	s.SetRegion("svc", "us-west-2")

	s2, _ := New(path)
	region, err := s2.GetRegion("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if region != "us-west-2" {
		t.Errorf("expected us-west-2, got %q", region)
	}
}
