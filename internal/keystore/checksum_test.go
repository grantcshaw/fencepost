package keystore_test

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/corecheck/fencepost/internal/keystore"
)

func TestSetChecksum_StoresChecksum(t *testing.T) {
	store := newTestStore(t)
	store.Set("svc", "key123")

	if err := store.SetChecksum("svc", "abc123"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := store.GetChecksum("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "abc123" {
		t.Errorf("expected abc123, got %s", got)
	}
}

func TestSetChecksum_MissingService(t *testing.T) {
	store := newTestStore(t)
	if err := store.SetChecksum("ghost", "abc"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestGetChecksum_MissingService(t *testing.T) {
	store := newTestStore(t)
	_, err := store.GetChecksum("ghost")
	if err == nil {
		t.Error("expected error for missing service")
	}
}

func TestClearChecksum_RemovesChecksum(t *testing.T) {
	store := newTestStore(t)
	store.Set("svc", "key123")
	store.SetChecksum("svc", "abc123")

	if err := store.ClearChecksum("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ := store.GetChecksum("svc")
	if got != "" {
		t.Errorf("expected empty checksum after clear, got %s", got)
	}
}

func TestComputeChecksum_ProducesCorrectDigest(t *testing.T) {
	store := newTestStore(t)
	store.Set("svc", "supersecretkey")

	got, err := store.ComputeChecksum("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sum := sha256.Sum256([]byte("supersecretkey"))
	want := fmt.Sprintf("%x", sum)
	if got != want {
		t.Errorf("expected %s, got %s", want, got)
	}
	stored, _ := store.GetChecksum("svc")
	if stored != want {
		t.Errorf("stored checksum mismatch: expected %s, got %s", want, stored)
	}
}

func TestServicesWithChecksum_ReturnsMatchingSorted(t *testing.T) {
	store := newTestStore(t)
	store.Set("bravo", "k1")
	store.Set("alpha", "k2")
	store.Set("charlie", "k3")
	store.SetChecksum("bravo", "sum1")
	store.SetChecksum("charlie", "sum2")

	result := store.ServicesWithChecksum()
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
	if result[0] != "bravo" || result[1] != "charlie" {
		t.Errorf("unexpected order: %v", result)
	}
}

func newTestStore(t *testing.T) *keystore.Store {
	t.Helper()
	path := tempStorePath(t)
	store, err := keystore.New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	return store
}
