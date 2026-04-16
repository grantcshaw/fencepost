package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func writeTempStore(t *testing.T, entries map[string]string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "store.json")
	ks := keystore.New(p)
	for k, v := range entries {
		if err := ks.Set(k, v); err != nil {
			t.Fatalf("set %s: %v", k, err)
		}
	}
	return p
}

func TestMergeCmd_AddsServices(t *testing.T) {
	srcPath := writeTempStore(t, map[string]string{"svc-a": "key-a", "svc-b": "key-b"})
	dstPath := writeTempStore(t, map[string]string{})

	cfgPath := writeTempConfig(t, dstPath)
	t.Setenv("FENCEPOST_CONFIG", cfgPath)

	Execute([]string{"merge", srcPath})

	ks := keystore.New(dstPath)
	e, err := ks.Get("svc-a")
	if err != nil {
		t.Fatalf("svc-a not found after merge: %v", err)
	}
	if e.Key != "key-a" {
		t.Errorf("expected key-a, got %s", e.Key)
	}
}

func TestMergeCmd_SkipsWithoutOverwrite(t *testing.T) {
	srcPath := writeTempStore(t, map[string]string{"svc-a": "new-key"})
	dstPath := writeTempStore(t, map[string]string{"svc-a": "original"})

	cfgPath := writeTempConfig(t, dstPath)
	t.Setenv("FENCEPOST_CONFIG", cfgPath)

	Execute([]string{"merge", srcPath})

	ks := keystore.New(dstPath)
	e, _ := ks.Get("svc-a")
	if e.Key != "original" {
		t.Errorf("key should not have changed without --overwrite")
	}
}

func writeTempConfig(t *testing.T, storePath string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "config.json")
	data, _ := json.Marshal(map[string]string{
		"store_path":     storePath,
		"audit_log_path": filepath.Join(t.TempDir(), "audit.log"),
	})
	_ = os.WriteFile(p, data, 0600)
	return p
}
