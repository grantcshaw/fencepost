package keystore_test

import (
	"testing"

	"github.com/clikd-inc/fencepost/internal/keystore"
)

func TestSetTier_StoresTier(t *testing.T) {
	ks := keystore.New(tempStorePath(t))
	ks.Set("svcA", "key1")
	if err := ks.SetTier("svcA", keystore.TierPro); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tier, err := ks.GetTier("svcA")
	if err != nil || tier != keystore.TierPro {
		t.Fatalf("expected pro, got %v (err=%v)", tier, err)
	}
}

func TestSetTier_MissingService(t *testing.T) {
	ks := keystore.New(tempStorePath(t))
	if err := ks.SetTier("ghost", keystore.TierBasic); err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestSetTier_InvalidValue(t *testing.T) {
	ks := keystore.New(tempStorePath(t))
	ks.Set("svcA", "key1")
	if err := ks.SetTier("svcA", keystore.Tier("platinum")); err == nil {
		t.Fatal("expected error for invalid tier")
	}
}

func TestGetTier_MissingService(t *testing.T) {
	ks := keystore.New(tempStorePath(t))
	_, err := ks.GetTier("ghost")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestClearTier_RemovesTier(t *testing.T) {
	ks := keystore.New(tempStorePath(t))
	ks.Set("svcA", "key1")
	ks.SetTier("svcA", keystore.TierEnterprise)
	if err := ks.ClearTier("svcA"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tier, _ := ks.GetTier("svcA")
	if tier != "" {
		t.Fatalf("expected empty tier, got %q", tier)
	}
}

func TestServicesByTier_ReturnsMatchingSorted(t *testing.T) {
	ks := keystore.New(tempStorePath(t))
	ks.Set("charlie", "k1")
	ks.Set("alpha", "k2")
	ks.Set("bravo", "k3")
	ks.SetTier("charlie", keystore.TierFree)
	ks.SetTier("alpha", keystore.TierFree)
	ks.SetTier("bravo", keystore.TierPro)
	result := ks.ServicesByTier(keystore.TierFree)
	if len(result) != 2 || result[0] != "alpha" || result[1] != "charlie" {
		t.Fatalf("unexpected result: %v", result)
	}
}
