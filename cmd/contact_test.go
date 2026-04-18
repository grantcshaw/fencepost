package cmd_test

import (
	"testing"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func newContactTestStore(t *testing.T) (string, string) {
	t.Helper()
	storePath := writeTempStore(t)
	cfgPath := writeTempConfig(t, storePath)
	return storePath, cfgPath
}

func TestContactCmd_SetAndGet(t *testing.T) {
	storePath, _ := newContactTestStore(t)
	ks := keystore.New(storePath)
	ks.Set("mysvc", "apikey")

	out, err := runCmd(t, "contact", "set", "mysvc", "dev@example.com")
	if err != nil {
		t.Fatalf("set failed: %v\n%s", err, out)
	}

	out, err = runCmd(t, "contact", "get", "mysvc")
	if err != nil {
		t.Fatalf("get failed: %v\n%s", err, out)
	}
	if out != "dev@example.com\n" {
		t.Errorf("expected dev@example.com, got %q", out)
	}
}

func TestContactCmd_ListByContact(t *testing.T) {
	storePath, _ := newContactTestStore(t)
	ks := keystore.New(storePath)
	ks.Set("alpha", "k1")
	ks.Set("beta", "k2")
	ks.SetContact("alpha", "team@corp.com")
	ks.SetContact("beta", "team@corp.com")

	out, err := runCmd(t, "contact", "list", "team@corp.com")
	if err != nil {
		t.Fatalf("list failed: %v\n%s", err, out)
	}
	if out != "alpha\nbeta\n" {
		t.Errorf("unexpected output: %q", out)
	}
}
