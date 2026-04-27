package cmd_test

import (
	"testing"
	"time"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func newReviewedTestStore(t *testing.T) *keystore.Store {
	t.Helper()
	s, err := keystore.New(t.TempDir() + "/store.json")
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	return s
}

func TestReviewedCmd_SetRecordsTimestamp(t *testing.T) {
	s := newReviewedTestStore(t)
	if err := s.Set("mysvc", "apikey"); err != nil {
		t.Fatalf("set failed: %v", err)
	}
	before := time.Now().Add(-time.Second)
	if err := s.SetReviewedAt("mysvc", time.Now()); err != nil {
		t.Fatalf("SetReviewedAt failed: %v", err)
	}
	got, err := s.GetReviewedAt("mysvc")
	if err != nil {
		t.Fatalf("GetReviewedAt failed: %v", err)
	}
	if got.Before(before) {
		t.Errorf("timestamp %v is before expected lower bound %v", got, before)
	}
}

func TestReviewedCmd_NeverReviewedLists(t *testing.T) {
	s := newReviewedTestStore(t)
	s.Set("alpha", "k1")
	s.Set("beta", "k2")
	s.SetReviewedAt("alpha", time.Now())
	result := s.NeverReviewed()
	if len(result) != 1 || result[0] != "beta" {
		t.Errorf("expected [beta], got %v", result)
	}
}

func TestReviewedCmd_ClearRemovesTimestamp(t *testing.T) {
	s := newReviewedTestStore(t)
	s.Set("svc", "key")
	s.SetReviewedAt("svc", time.Now())
	if err := s.ClearReviewedAt("svc"); err != nil {
		t.Fatalf("ClearReviewedAt failed: %v", err)
	}
	got, _ := s.GetReviewedAt("svc")
	if !got.IsZero() {
		t.Errorf("expected zero time, got %v", got)
	}
}
