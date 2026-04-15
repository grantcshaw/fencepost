package keystore

import (
	"testing"
	"time"
)

func TestDiff_AddedService(t *testing.T) {
	s := newTestStore(t)
	base := s.SnapshotData()

	_ = s.Set("newservice", "key-new")

	result := s.Diff(base)
	if len(result.Added) != 1 || result.Added[0] != "newservice" {
		t.Errorf("expected Added=[newservice], got %v", result.Added)
	}
	if len(result.Removed)+len(result.Rotated)+len(result.Modified) != 0 {
		t.Errorf("unexpected changes: %+v", result)
	}
}

func TestDiff_RemovedService(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key-abc")
	base := s.SnapshotData()

	delete(s.data, "svc")

	result := s.Diff(base)
	if len(result.Removed) != 1 || result.Removed[0] != "svc" {
		t.Errorf("expected Removed=[svc], got %v", result.Removed)
	}
}

func TestDiff_RotatedService(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key-old")
	base := s.SnapshotData()

	// Simulate a rotation
	entry := s.data["svc"]
	entry.Key = "key-new"
	entry.RotatedAt = time.Now()
	s.data["svc"] = entry

	result := s.Diff(base)
	if len(result.Rotated) != 1 || result.Rotated[0] != "svc" {
		t.Errorf("expected Rotated=[svc], got %v", result.Rotated)
	}
	if len(result.Modified) != 0 {
		t.Errorf("expected no Modified, got %v", result.Modified)
	}
}

func TestDiff_ModifiedService(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key-old")
	base := s.SnapshotData()

	entry := s.data["svc"]
	entry.Key = "key-changed"
	s.data["svc"] = entry

	result := s.Diff(base)
	if len(result.Modified) != 1 || result.Modified[0] != "svc" {
		t.Errorf("expected Modified=[svc], got %v", result.Modified)
	}
}

func TestDiff_NoChanges(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key-abc")
	base := s.SnapshotData()

	result := s.Diff(base)
	if len(result.Added)+len(result.Removed)+len(result.Rotated)+len(result.Modified) != 0 {
		t.Errorf("expected no changes, got %+v", result)
	}
}

func newTestStore(t *testing.T) *Store {
	t.Helper()
	s, err := New(tempStorePath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}
