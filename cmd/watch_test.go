package cmd

import (
	"bytes"
	"testing"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func TestWatchCmd_AddsService(t *testing.T) {
	store, cfg := writeTempStore(t)
	_ = store.Set("myapi", "secret")

	rootCmd.SetArgs([]string{"--config", cfg, "watch", "myapi"})
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("watch cmd: %v", err)
	}

	s2, _ := keystore.New(store.Path())
	if !s2.IsWatched("myapi") {
		t.Error("expected myapi to be watched")
	}
}

func TestUnwatchCmd_RemovesService(t *testing.T) {
	store, cfg := writeTempStore(t)
	_ = store.Set("myapi", "secret")
	_ = store.Watch("myapi")

	rootCmd.SetArgs([]string{"--config", cfg, "unwatch", "myapi"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unwatch cmd: %v", err)
	}

	s2, _ := keystore.New(store.Path())
	if s2.IsWatched("myapi") {
		t.Error("expected myapi to not be watched")
	}
}
