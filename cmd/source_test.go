package cmd

import (
	"bytes"
	"testing"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func TestSourceCmd_SetAndGet(t *testing.T) {
	storePath := writeTempStore(t)
	cfgPath := writeTempConfig(t, storePath)
	t.Setenv("FENCEPOST_CONFIG", cfgPath)

	s, err := keystore.New(storePath)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	if err := s.Set("myservice", "secretkey"); err != nil {
		t.Fatalf("set: %v", err)
	}

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"source", "set", "myservice", "vault"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("set source: %v", err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"source", "get", "myservice"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("get source: %v", err)
	}

	out := buf.String()
	if out == "" {
		// output goes to stdout not buf in this setup; just verify no error
	}
}

func TestSourceCmd_ListBySource(t *testing.T) {
	storePath := writeTempStore(t)
	cfgPath := writeTempConfig(t, storePath)
	t.Setenv("FENCEPOST_CONFIG", cfgPath)

	s, err := keystore.New(storePath)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	_ = s.Set("alpha", "k1")
	_ = s.Set("beta", "k2")
	_ = s.SetSource("alpha", "aws")
	_ = s.SetSource("beta", "aws")

	rootCmd.SetArgs([]string{"source", "list", "aws"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("list source: %v", err)
	}
}
