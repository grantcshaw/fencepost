package keystore

import (
	"testing"
	"time"
)

func TestSetLastAccessed_UpdatesTimestamp(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key123")

	before := time.Now().UTC().Add(-time.Second)
	if err := s.SetLastAccessed("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ts, err := s.GetLastAccessed("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ts.After(before) {
		t.Errorf("expected timestamp after %v, got %v", before, ts)
	}
}

func TestSetLastAccessed_MissingService(t *testing.T) {
	s := newTestStore(t)
	if err := s.SetLastAccessed("ghost"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestGetLastAccessed_MissingService(t *testing.T) {
	s := newTestStore(t)
	_, err := s.GetLastAccessed("ghost")
	if err == nil {
		t.Error("expected error for missing service")
	}
}

func TestNeverAccessed_ReturnsUnaccessed(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("alpha", "k1")
	_ = s.Set("beta", "k2")
	_ = s.SetLastAccessed("alpha")

	result := s.NeverAccessed()
	if len(result) != 1 || result[0] != "beta" {
		t.Errorf("expected [beta], got %v", result)
	}
}

func TestAccessedSince_ReturnsMatchingServices(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc1", "k1")
	_ = s.Set("svc2", "k2")

	cutoff := time.Now().UTC()
	time.Sleep(10 * time.Millisecond)
	_ = s.SetLastAccessed("svc1")

	result := s.AccessedSince(cutoff)
	if len(result) != 1 || result[0] != "svc1" {
		t.Errorf("expected [svc1], got %v", result)
	}
}

func TestNeverAccessed_PersistsAcrossReload(t *testing.T) {
	path := tempStorePath(t)
	s, _ := New(path)
	_ = s.Set("svc", "key")

	s2, _ := New(path)
	result := s2.NeverAccessed()
	if len(result) != 1 || result[0] != "svc" {
		t.Errorf("expected [svc], got %v", result)
	}
}
