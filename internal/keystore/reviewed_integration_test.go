package keystore_test

import (
	"testing"
	"time"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func TestReviewed_PersistsToDisk(t *testing.T) {
	path := t.TempDir() + "/store.json"
	s, _ := keystore.New(path)
	s.Set("svc", "key")
	now := time.Now().Truncate(time.Second)
	s.SetReviewedAt("svc", now)

	s2, err := keystore.New(path)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	got, err := s2.GetReviewedAt("svc")
	if err != nil {
		t.Fatalf("GetReviewedAt after reload: %v", err)
	}
	if !got.Equal(now) {
		t.Errorf("expected %v after reload, got %v", now, got)
	}
}

func TestReviewed_ClearPersistsToDisk(t *testing.T) {
	path := t.TempDir() + "/store.json"
	s, _ := keystore.New(path)
	s.Set("svc", "key")
	s.SetReviewedAt("svc", time.Now())
	s.ClearReviewedAt("svc")

	s2, _ := keystore.New(path)
	got, _ := s2.GetReviewedAt("svc")
	if !got.IsZero() {
		t.Errorf("expected zero time after reload, got %v", got)
	}
}

func TestReviewedBefore_PersistsAndFilters(t *testing.T) {
	path := t.TempDir() + "/store.json"
	s, _ := keystore.New(path)
	s.Set("stale", "k1")
	s.Set("fresh", "k2")
	s.SetReviewedAt("stale", time.Now().Add(-60*24*time.Hour))
	s.SetReviewedAt("fresh", time.Now())

	s2, _ := keystore.New(path)
	cutoff := time.Now().Add(-7 * 24 * time.Hour)
	result := s2.ReviewedBefore(cutoff)
	if len(result) != 1 || result[0] != "stale" {
		t.Errorf("expected [stale], got %v", result)
	}
}
