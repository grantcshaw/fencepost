package audit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/fencepost/internal/audit"
)

func tempLogPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "audit.log")
}

func TestNew_CreatesFile(t *testing.T) {
	path := tempLogPath(t)
	_, err := audit.New(path)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected log file to be created at %s", path)
	}
}

func TestLogger_Log_WritesEntry(t *testing.T) {
	path := tempLogPath(t)
	logger, err := audit.New(path)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	err = logger.Log(audit.EventKeyCreated, "stripe", "key_abc123", "created new API key")
	if err != nil {
		t.Fatalf("Log() error = %v", err)
	}

	entries, err := audit.ReadAll(path)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Event != audit.EventKeyCreated {
		t.Errorf("expected event %q, got %q", audit.EventKeyCreated, entries[0].Event)
	}
	if entries[0].Service != "stripe" {
		t.Errorf("expected service %q, got %q", "stripe", entries[0].Service)
	}
}

func TestReadAll_EmptyFile(t *testing.T) {
	path := tempLogPath(t)
	entries, err := audit.ReadAll(path)
	if err != nil {
		t.Fatalf("ReadAll() on missing file error = %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestFilterByService(t *testing.T) {
	path := tempLogPath(t)
	logger, _ := audit.New(path)
	_ = logger.Log(audit.EventKeyCreated, "stripe", "k1", "msg")
	_ = logger.Log(audit.EventKeyRotated, "github", "k2", "msg")
	_ = logger.Log(audit.EventKeyRevoked, "stripe", "k3", "msg")

	entries, _ := audit.ReadAll(path)
	filtered := audit.FilterByService(entries, "stripe")
	if len(filtered) != 2 {
		t.Errorf("expected 2 stripe entries, got %d", len(filtered))
	}
}
