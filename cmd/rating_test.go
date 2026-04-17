package cmd_test

import (
	"testing"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func TestRatingCmd_SetAndGet(t *testing.T) {
	store := keystore.New(tempStorePath(t))
	store.Set("api-gateway", "secretkey")

	if err := store.SetRating("api-gateway", "critical"); err != nil {
		t.Fatalf("SetRating failed: %v", err)
	}

	r, err := store.GetRating("api-gateway")
	if err != nil {
		t.Fatalf("GetRating failed: %v", err)
	}
	if r != "critical" {
		t.Errorf("expected critical, got %s", r)
	}
}

func TestRatingCmd_ListByRating(t *testing.T) {
	store := keystore.New(tempStorePath(t))
	store.Set("svc-a", "k1")
	store.Set("svc-b", "k2")
	store.Set("svc-c", "k3")
	store.SetRating("svc-a", "low")
	store.SetRating("svc-c", "low")

	results, err := store.ServicesByRating("low")
	if err != nil {
		t.Fatalf("ServicesByRating failed: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
	if results[0] != "svc-a" || results[1] != "svc-c" {
		t.Errorf("unexpected order: %v", results)
	}
}
