package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func runAccessCmd(args ...string) (string, error) {
	buf := &bytes.Buffer{}
	rootCmd := &cobra.Command{Use: "fencepost"}
	rootCmd.AddCommand(accessCmd)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	return buf.String(), err
}

func TestAccessTouch_RecordsTimestamp(t *testing.T) {
	cfgPath, storePath := writeTempAccessConfig(t)
	_ = cfgPath
	_ = storePath

	// seed a service
	ks := newAccessTestStore(t, storePath)
	_ = ks.Set("mysvc", "abc123")

	t.Setenv("FENCEPOST_CONFIG", cfgPath)
	_, err := runAccessCmd("access", "touch", "mysvc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ts, err := ks.GetLastAccessed("mysvc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ts.IsZero() {
		t.Error("expected non-zero last accessed timestamp")
	}
}

func TestAccessNever_ListsUnaccessed(t *testing.T) {
	cfgPath, storePath := writeTempAccessConfig(t)
	_ = storePath

	ks := newAccessTestStore(t, storePath)
	_ = ks.Set("alpha", "k1")
	_ = ks.Set("beta", "k2")
	_ = ks.SetLastAccessed("alpha")

	t.Setenv("FENCEPOST_CONFIG", cfgPath)
	out, err := runAccessCmd("access", "never")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "beta") {
		t.Errorf("expected beta in output, got: %s", out)
	}
	if strings.Contains(out, "alpha") {
		t.Errorf("did not expect alpha in output, got: %s", out)
	}
}
