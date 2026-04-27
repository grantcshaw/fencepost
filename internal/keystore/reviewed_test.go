package keystore_test

import (
	"testing"
	"time"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func TestSetReviewedAt_StoresTimestamp(t *testing.T) {
	s := newTestStore(t)
	s.Set("svc", "key123")
	now := time.Now().Truncate(time.Second)
	if err := s.SetReviewedAt("svc", now); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := s.GetReviewedAt("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.Equal(now) {
		t.Errorf("expected %v, got %v", now, got)
	}
}

func TestSetReviewedAt_MissingService(t *testing.T) {
	s := newTestStore(t)
	err := s.SetReviewedAt("ghost", time.Now())
	if err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestGetReviewedAt_MissingService(t *testing.T) {
	s := newTestStore(t)
	_, err := s.GetReviewedAt("ghost")
	if err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestClearReviewedAt_RemovesTimestamp(t *testing.T) {
	s := newTestStore(t)
	s.Set("svc", "key123")
	s.SetReviewedAt("svc", time.Now())
	if err := s.ClearReviewedAt("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ := s.GetReviewedAt("svc")
	if !got.IsZero() {
		t.Errorf("expected zero time after clear, got %v", got)
	}
}

func TestNeverReviewed_ReturnsUnreviewed(t *testing.T) {
	s := newTestStore(t)
	s.Set("alpha", "k1")
	s.Set("beta", "k2")
	s.Set("gamma", "k3")
	s.SetReviewedAt("beta", time.Now())
	result := s.NeverReviewed()
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d: %v", len(result), result)
	}
	if result[0] != "alpha" || result[1] != "gamma" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestReviewedBefore_ReturnsStale(t *testing.T) {
	s := newTestStore(t)
	s.Set("old", "k1")
	s.Set("new", "k2")
	past := time.Now().Add(-30 * 24 * time.Hour)
	s.SetReviewedAt("old", past)
	s.SetReviewedAt("new", time.Now())
	cutoff := time.Now().Add(-7 * 24 * time.Hour)
	result := s.ReviewedBefore(cutoff)
	if len(result) != 1 || result[0] != "old" {
		t.Errorf("expected [old], got %v", result)
	}
}
