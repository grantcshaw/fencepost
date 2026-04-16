package keystore

import (
	"testing"
)

func TestSetGroup_StoresGroup(t *testing.T) {
	s, _ := New(tempStorePath(t))
	s.Set("svcA", "key1")

	if err := s.SetGroup("svcA", "payments"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	g, _ := s.GetGroup("svcA")
	if g != "payments" {
		t.Errorf("expected 'payments', got %q", g)
	}
}

func TestSetGroup_MissingService(t *testing.T) {
	s, _ := New(tempStorePath(t))
	if err := s.SetGroup("ghost", "infra"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestGetGroup_MissingService(t *testing.T) {
	s, _ := New(tempStorePath(t))
	_, err := s.GetGroup("ghost")
	if err == nil {
		t.Error("expected error for missing service")
	}
}

func TestClearGroup_RemovesGroup(t *testing.T) {
	s, _ := New(tempStorePath(t))
	s.Set("svcA", "key1")
	s.SetGroup("svcA", "payments")

	if err := s.ClearGroup("svcA"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	g, _ := s.GetGroup("svcA")
	if g != "" {
		t.Errorf("expected empty group, got %q", g)
	}
}

func TestServicesByGroup_ReturnsMatchingSorted(t *testing.T) {
	s, _ := New(tempStorePath(t))
	s.Set("svcB", "k2")
	s.Set("svcA", "k1")
	s.Set("svcC", "k3")
	s.SetGroup("svcB", "infra")
	s.SetGroup("svcA", "infra")
	s.SetGroup("svcC", "payments")

	result := s.ServicesByGroup("infra")
	if len(result) != 2 || result[0] != "svcA" || result[1] != "svcB" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestGroups_ReturnsDistinctSorted(t *testing.T) {
	s, _ := New(tempStorePath(t))
	s.Set("svcA", "k1")
	s.Set("svcB", "k2")
	s.Set("svcC", "k3")
	s.SetGroup("svcA", "infra")
	s.SetGroup("svcB", "payments")
	s.SetGroup("svcC", "infra")

	groups := s.Groups()
	if len(groups) != 2 || groups[0] != "infra" || groups[1] != "payments" {
		t.Errorf("unexpected groups: %v", groups)
	}
}
