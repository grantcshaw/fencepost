package cmd_test

import (
	"testing"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func newRevokedTestStore(t *testing.T) *keystore.Store {
	t.Helper()
	path := tempStorePath(t)
	s, err := keystore.New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	return s
}

func TestRevokedCmd_SetAndGet(t *testing.T) {
	s := newRevokedTestStore(t)
	if err := s.Set("myapi", "secretkey"); err != nil {
		t.Fatalf("set failed: %v", err)
	}

	if err := s.SetRevoked("myapi", "key was leaked"); err != nil {
		t.Fatalf("SetRevoked failed: %v", err)
	}

	ok, err := s.IsRevoked("myapi")
	if err != nil {
		t.Fatalf("IsRevoked failed: %v", err)
	}
	if !ok {
		t.Error("expected service to be revoked")
	}
}

func TestRevokedCmd_ClearRevoked(t *testing.T) {
	s := newRevokedTestStore(t)
	_ = s.Set("myapi", "secretkey")
	_ = s.SetRevoked("myapi", "test")

	if err := s.ClearRevoked("myapi"); err != nil {
		t.Fatalf("ClearRevoked failed: %v", err)
	}

	ok, _ := s.IsRevoked("myapi")
	if ok {
		t.Error("expected revoked to be cleared")
	}
}

func TestRevokedCmd_ListRevoked(t *testing.T) {
	s := newRevokedTestStore(t)
	for _, svc := range []string{"svc1", "svc2", "svc3"} {
		_ = s.Set(svc, "k")
		_ = s.SetRevoked(svc, "batch revoke")
	}
	_ = s.Set("svc4", "k") // not revoked

	keys := s.RevokedKeys()
	if len(keys) != 3 {
		t.Fatalf("expected 3 revoked keys, got %d", len(keys))
	}
}
