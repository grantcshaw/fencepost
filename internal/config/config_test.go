package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fencepost/internal/config"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}
	return path
}

func TestLoad_ValidConfig(t *testing.T) {
	cfgContent := `
audit_log_path: /tmp/audit.log
services:
  github:
    name: GitHub
    key_file: /tmp/github.key
    rotate_days: 30
`
	path := writeTempConfig(t, cfgContent)

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.AuditLogPath != "/tmp/audit.log" {
		t.Errorf("expected audit_log_path /tmp/audit.log, got %s", cfg.AuditLogPath)
	}

	svc, ok := cfg.Services["github"]
	if !ok {
		t.Fatal("expected github service in config")
	}
	if svc.RotateDays != 30 {
		t.Errorf("expected rotate_days 30, got %d", svc.RotateDays)
	}
}

func TestLoad_DefaultAuditLogPath(t *testing.T) {
	cfgContent := `services: {}`
	path := writeTempConfig(t, cfgContent)

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.AuditLogPath == "" {
		t.Error("expected a default audit_log_path, got empty string")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("expected error for missing config file, got nil")
	}
}
