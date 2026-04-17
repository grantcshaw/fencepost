package keystore

import (
	"testing"
)

func TestSetSource_StoresSource(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svcA", "key1")
	if err := s.SetSource("svcA", "vault"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	src, err := s.GetSource("svcA")
	if err != nil || src != "vault" {
		t.Fatalf("expected vault, got %q err %v", src, err)
	}
}

func TestSetSource_MissingService(t *testing.T) {
	s := newTestStore(t)
	if err := s.SetSource("missing", "aws"); err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestSetSource_InvalidValue(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svcA", "key1")
	if err := s.SetSource("svcA", "magic"); err == nil {
		t.Fatal("expected error for invalid source")
	}
}

func TestGetSource_DefaultsToManual(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svcA", "key1")
	src, err := s.GetSource("svcA")
	if err != nil || src != "manual" {
		t.Fatalf("expected manual default, got %q err %v", src, err)
	}
}

func TestGetSource_MissingService(t *testing.T) {
	s := newTestStore(t)
	_, err := s.GetSource("ghost")
	if err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestServicesBySource_ReturnsMatchingSorted(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svcB", "k")
	_ = s.Set("svcA", "k")
	_ = s.Set("svcC", "k")
	_ = s.SetSource("svcA", "aws")
	_ = s.SetSource("svcC", "aws")
	results, err := s.ServicesBySource("aws")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 || results[0] != "svcA" || results[1] != "svcC" {
		t.Fatalf("unexpected results: %v", results)
	}
}

func TestClearSource_RemovesSource(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svcA", "key1")
	_ = s.SetSource("svcA", "gcp")
	_ = s.ClearSource("svcA")
	src, _ := s.GetSource("svcA")
	if src != "manual" {
		t.Fatalf("expected manual after clear, got %q", src)
	}
}
