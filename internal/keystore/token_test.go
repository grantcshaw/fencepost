package keystore_test

import (
	"testing"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func TestSetToken_StoresToken(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	s.Set("svc", "key123")

	if err := s.SetToken("svc", "tok-abc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := s.GetToken("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "tok-abc" {
		t.Errorf("expected tok-abc, got %s", got)
	}
}

func TestSetToken_MissingService(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	if err := s.SetToken("ghost", "tok"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestGetToken_MissingService(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	if _, err := s.GetToken("ghost"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestClearToken_RemovesToken(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	s.Set("svc", "key123")
	s.SetToken("svc", "tok-abc")

	if err := s.ClearToken("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ := s.GetToken("svc")
	if got != "" {
		t.Errorf("expected empty token, got %s", got)
	}
}

func TestServicesWithToken_ReturnsMatchingSorted(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	s.Set("bravo", "k1")
	s.Set("alpha", "k2")
	s.Set("charlie", "k3")
	s.SetToken("bravo", "t1")
	s.SetToken("charlie", "t2")

	result := s.ServicesWithToken()
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
	if result[0] != "bravo" || result[1] != "charlie" {
		t.Errorf("unexpected order: %v", result)
	}
}
