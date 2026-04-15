package audit_test

import (
	"testing"
	"time"

	"github.com/user/fencepost/internal/audit"
)

func TestReadAll_MultipleEntries(t *testing.T) {
	path := tempLogPath(t)
	logger, err := audit.New(path)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	if err := logger.Log("svc-a", "key.created", "created key for svc-a"); err != nil {
		t.Fatalf("Log() error: %v", err)
	}
	if err := logger.Log("svc-b", "key.rotated", "rotated key for svc-b"); err != nil {
		t.Fatalf("Log() error: %v", err)
	}

	entries, err := audit.ReadAll(path)
	if err != nil {
		t.Fatalf("ReadAll() error: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

func TestFilterByService_ReturnsMatchingEntries(t *testing.T) {
	path := tempLogPath(t)
	logger, err := audit.New(path)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	_ = logger.Log("svc-a", "key.created", "msg")
	_ = logger.Log("svc-b", "key.rotated", "msg")
	_ = logger.Log("svc-a", "key.rotated", "msg")

	entries, err := audit.ReadAll(path)
	if err != nil {
		t.Fatalf("ReadAll() error: %v", err)
	}

	filtered := audit.FilterByService(entries, "svc-a")
	if len(filtered) != 2 {
		t.Errorf("expected 2 entries for svc-a, got %d", len(filtered))
	}
	for _, e := range filtered {
		if e.Service != "svc-a" {
			t.Errorf("unexpected service %q in filtered results", e.Service)
		}
	}
}

func TestFilterByEvent_ReturnsMatchingEntries(t *testing.T) {
	path := tempLogPath(t)
	logger, err := audit.New(path)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	_ = logger.Log("svc-a", "key.created", "msg")
	_ = logger.Log("svc-b", "key.rotated", "msg")
	_ = logger.Log("svc-c", "key.rotated", "msg")

	entries, err := audit.ReadAll(path)
	if err != nil {
		t.Fatalf("ReadAll() error: %v", err)
	}

	filtered := audit.FilterByEvent(entries, "key.rotated")
	if len(filtered) != 2 {
		t.Errorf("expected 2 entries for key.rotated, got %d", len(filtered))
	}
}

func TestReadAll_TimestampsAreParsed(t *testing.T) {
	path := tempLogPath(t)
	logger, err := audit.New(path)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	before := time.Now().Truncate(time.Second)
	_ = logger.Log("svc-a", "key.created", "msg")
	after := time.Now().Add(time.Second)

	entries, err := audit.ReadAll(path)
	if err != nil {
		t.Fatalf("ReadAll() error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	ts := entries[0].Timestamp
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v out of expected range [%v, %v]", ts, before, after)
	}
}
