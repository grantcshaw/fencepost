package keystore_test

import (
	"testing"

	"github.com/danielmichaels/fencepost/internal/keystore"
)

func TestSetFlag_StoresFlag(t *testing.T) {
	ks, _ := keystore.New(tempStorePath(t))
	_ = ks.Set("svc", "key123")
	if err := ks.SetFlag("svc", "critical"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	flags, _ := ks.GetFlags("svc")
	if len(flags) != 1 || flags[0] != "critical" {
		t.Errorf("expected [critical], got %v", flags)
	}
}

func TestSetFlag_MissingService(t *testing.T) {
	ks, _ := keystore.New(tempStorePath(t))
	if err := ks.SetFlag("ghost", "critical"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestSetFlag_NoDuplicates(t *testing.T) {
	ks, _ := keystore.New(tempStorePath(t))
	_ = ks.Set("svc", "key123")
	_ = ks.SetFlag("svc", "critical")
	_ = ks.SetFlag("svc", "critical")
	flags, _ := ks.GetFlags("svc")
	if len(flags) != 1 {
		t.Errorf("expected 1 flag, got %d", len(flags))
	}
}

func TestUnsetFlag_RemovesFlag(t *testing.T) {
	ks, _ := keystore.New(tempStorePath(t))
	_ = ks.Set("svc", "key123")
	_ = ks.SetFlag("svc", "critical")
	_ = ks.UnsetFlag("svc", "critical")
	flags, _ := ks.GetFlags("svc")
	if len(flags) != 0 {
		t.Errorf("expected no flags, got %v", flags)
	}
}

func TestGetFlags_MissingService(t *testing.T) {
	ks, _ := keystore.New(tempStorePath(t))
	_, err := ks.GetFlags("ghost")
	if err == nil {
		t.Error("expected error for missing service")
	}
}

func TestServicesByFlag_ReturnsMatchingSorted(t *testing.T) {
	ks, _ := keystore.New(tempStorePath(t))
	_ = ks.Set("svc-b", "k1")
	_ = ks.Set("svc-a", "k2")
	_ = ks.Set("svc-c", "k3")
	_ = ks.SetFlag("svc-b", "urgent")
	_ = ks.SetFlag("svc-a", "urgent")
	result := ks.ServicesByFlag("urgent")
	if len(result) != 2 || result[0] != "svc-a" || result[1] != "svc-b" {
		t.Errorf("unexpected result: %v", result)
	}
}
