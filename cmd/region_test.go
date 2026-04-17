package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fencepost/internal/keystore"
)

func TestRegionCmd_SetAndGet(t *testing.T) {
	dir := t.TempDir()
	storePath := filepath.Join(dir, "store.json")
	cfgPath := writeTempConfig(t, dir, storePath)

	s, _ := keystore.New(storePath)
	s.Set("mysvc", "apikey")

	os.Args = []string{"fencepost", "--config", cfgPath, "region", "set", "mysvc", "eu-west-1"}
	if err := Execute(); err != nil {
		t.Fatalf("set region failed: %v", err)
	}

	s2, _ := keystore.New(storePath)
	region, err := s2.GetRegion("mysvc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if region != "eu-west-1" {
		t.Errorf("expected eu-west-1, got %q", region)
	}
}

func TestRegionCmd_ListByRegion(t *testing.T) {
	dir := t.TempDir()
	storePath := filepath.Join(dir, "store.json")
	cfgPath := writeTempConfig(t, dir, storePath)

	s, _ := keystore.New(storePath)
	s.Set("alpha", "k1")
	s.Set("beta", "k2")
	s.SetRegion("alpha", "us-east-1")
	s.SetRegion("beta", "us-east-1")

	os.Args = []string{"fencepost", "--config", cfgPath, "region", "list", "us-east-1"}
	if err := Execute(); err != nil {
		t.Fatalf("list region failed: %v", err)
	}
}
