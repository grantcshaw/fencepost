package cmd_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/clikd-inc/fencepost/internal/keystore"
)

func writeTempTierStore(t *testing.T) (storePath string, cfgPath string) {
	t.Helper()
	dir := t.TempDir()
	storePath = filepath.Join(dir, "store.json")
	cfgPath = filepath.Join(dir, "config.yaml")
	data := map[string]interface{}{
		"store_path": storePath,
		"audit_log_path": filepath.Join(dir, "audit.log"),
	}
	b, _ := json.Marshal(data)
	os.WriteFile(cfgPath, b, 0600)
	return
}

func TestTierCmd_SetAndGet(t *testing.T) {
	storePath, _ := writeTempTierStore(t)
	ks := keystore.New(storePath)
	ks.Set("myservice", "secretkey")

	if err := ks.SetTier("myservice", keystore.TierEnterprise); err != nil {
		t.Fatalf("SetTier failed: %v", err)
	}
	tier, err := ks.GetTier("myservice")
	if err != nil {
		t.Fatalf("GetTier failed: %v", err)
	}
	if tier != keystore.TierEnterprise {
		t.Fatalf("expected enterprise, got %q", tier)
	}
}

func TestTierCmd_ListByTier(t *testing.T) {
	storePath, _ := writeTempTierStore(t)
	ks := keystore.New(storePath)
	ks.Set("svc1", "k1")
	ks.Set("svc2", "k2")
	ks.Set("svc3", "k3")
	ks.SetTier("svc1", keystore.TierBasic)
	ks.SetTier("svc2", keystore.TierBasic)
	ks.SetTier("svc3", keystore.TierPro)

	result := ks.ServicesByTier(keystore.TierBasic)
	if len(result) != 2 {
		t.Fatalf("expected 2 basic services, got %d", len(result))
	}
}
