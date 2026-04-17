package cmd

import (
	"testing"

	"github.com/fencepost/internal/keystore"
)

func TestProtocolCmd_SetAndGet(t *testing.T) {
	storePath := writeTempStore(t)
	s, _ := keystore.New(storePath)
	s.Set("mysvc", "apikey")

	if err := s.SetProtocol("mysvc", "grpc"); err != nil {
		t.Fatalf("SetProtocol failed: %v", err)
	}

	p, err := s.GetProtocol("mysvc")
	if err != nil {
		t.Fatalf("GetProtocol failed: %v", err)
	}
	if p != "grpc" {
		t.Errorf("expected grpc, got %s", p)
	}
}

func TestProtocolCmd_ListByProtocol(t *testing.T) {
	storePath := writeTempStore(t)
	s, _ := keystore.New(storePath)
	s.Set("alpha", "k1")
	s.Set("beta", "k2")
	s.SetProtocol("alpha", "webhook")
	s.SetProtocol("beta", "grpc")

	results := s.ServicesByProtocol("webhook")
	if len(results) != 1 || results[0] != "alpha" {
		t.Errorf("unexpected results: %v", results)
	}
}
