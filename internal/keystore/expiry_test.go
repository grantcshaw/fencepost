package keystore

import (
	"testing"
	"time"
)

func TestIsExpired_FreshKey(t *testing.T) {
	policy := DefaultExpiryPolicy()
	entry := Entry{
		Key:       "fresh-key",
		CreatedAt: time.Now(),
	}
	if policy.IsExpired(entry) {
		t.Error("expected fresh key to not be expired")
	}
}

func TestIsExpired_OldKey(t *testing.T) {
	policy := DefaultExpiryPolicy()
	entry := Entry{
		Key:       "old-key",
		CreatedAt: time.Now().Add(-100 * 24 * time.Hour),
	}
	if !policy.IsExpired(entry) {
		t.Error("expected old key to be expired")
	}
}

func TestIsExpired_ZeroCreatedAt(t *testing.T) {
	policy := DefaultExpiryPolicy()
	entry := Entry{Key: "no-date"}
	if policy.IsExpired(entry) {
		t.Error("expected zero CreatedAt to never be expired")
	}
}

func TestExpiredKeys_ReturnsCorrectServices(t *testing.T) {
	path := tempStorePath(t)
	s, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	// fresh key
	if err := s.Set("svc-fresh", "key-fresh"); err != nil {
		t.Fatalf("Set: %v", err)
	}

	// manually insert an old entry
	s.mu.Lock()
	s.data["svc-old"] = Entry{
		Key:       "key-old",
		CreatedAt: time.Now().Add(-100 * 24 * time.Hour),
	}
	s.mu.Unlock()

	policy := DefaultExpiryPolicy()
	expired := ExpiredKeys(s, policy)

	if len(expired) != 1 {
		t.Fatalf("expected 1 expired key, got %d", len(expired))
	}
	if expired[0] != "svc-old" {
		t.Errorf("expected svc-old, got %s", expired[0])
	}
}

func TestExpiredKeys_EmptyStore(t *testing.T) {
	path := tempStorePath(t)
	s, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	expired := ExpiredKeys(s, DefaultExpiryPolicy())
	if len(expired) != 0 {
		t.Errorf("expected no expired keys, got %d", len(expired))
	}
}
