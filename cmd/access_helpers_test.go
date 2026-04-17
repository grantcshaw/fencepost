package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/cqroot/fencepost/internal/keystore"
)

func writeTempAccessConfig(t *testing.T) (cfgPath, storePath string) {
	t.Helper()
	dir := t.TempDir()
	storePath = filepath.Join(dir, "store.json")
	auditPath := filepath.Join(dir, "audit.log")
	cfgPath = filepath.Join(dir, "config.yaml")

	content := fmt.Sprintf("store_path: %s\naudit_log_path: %s\n", storePath, auditPath)
	if err := os.WriteFile(cfgPath, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	return cfgPath, storePath
}

func newAccessTestStore(t *testing.T, path string) *keystore.Store {
	t.Helper()
	s, err := keystore.New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	return s
}
