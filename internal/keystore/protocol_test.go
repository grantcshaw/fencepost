package keystore

import (
	"testing"
)

func TestSetProtocol_StoresProtocol(t *testing.T) {
	s := newTestStore(t)
	s.Set("svc", "key123")
	if err := s.SetProtocol("svc", "grpc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p, _ := s.GetProtocol("svc")
	if p != "grpc" {
		t.Errorf("expected grpc, got %s", p)
	}
}

func TestSetProtocol_MissingService(t *testing.T) {
	s := newTestStore(t)
	if err := s.SetProtocol("missing", "rest"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestSetProtocol_InvalidValue(t *testing.T) {
	s := newTestStore(t)
	s.Set("svc", "key123")
	if err := s.SetProtocol("svc", "ftp"); err == nil {
		t.Error("expected error for invalid protocol")
	}
}

func TestGetProtocol_DefaultsToRest(t *testing.T) {
	s := newTestStore(t)
	s.Set("svc", "key123")
	p, err := s.GetProtocol("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p != "rest" {
		t.Errorf("expected rest, got %s", p)
	}
}

func TestGetProtocol_MissingService(t *testing.T) {
	s := newTestStore(t)
	if _, err := s.GetProtocol("missing"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestServicesByProtocol_ReturnsMatchingSorted(t *testing.T) {
	s := newTestStore(t)
	s.Set("alpha", "k1")
	s.Set("beta", "k2")
	s.Set("gamma", "k3")
	s.SetProtocol("alpha", "grpc")
	s.SetProtocol("gamma", "grpc")
	results := s.ServicesByProtocol("grpc")
	if len(results) != 2 || results[0] != "alpha" || results[1] != "gamma" {
		t.Errorf("unexpected results: %v", results)
	}
}

func TestClearProtocol_RemovesProtocol(t *testing.T) {
	s := newTestStore(t)
	s.Set("svc", "key123")
	s.SetProtocol("svc", "graphql")
	s.ClearProtocol("svc")
	p, _ := s.GetProtocol("svc")
	if p != "rest" {
		t.Errorf("expected default rest after clear, got %s", p)
	}
}
