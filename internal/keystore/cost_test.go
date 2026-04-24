package keystore_test

import (
	"testing"

	"github.com/yourusername/fencepost/internal/keystore"
)

func TestSetCost_StoresCost(t *testing.T) {
	ks := keystore.New(tempStorePath(t))
	if err := ks.Set("svc", "key123"); err != nil {
		t.Fatal(err)
	}
	if err := ks.SetCost("svc", 4999); err != nil {
		t.Fatalf("SetCost: %v", err)
	}
	cost, err := ks.GetCost("svc")
	if err != nil {
		t.Fatalf("GetCost: %v", err)
	}
	if cost != 4999 {
		t.Errorf("expected 4999, got %d", cost)
	}
}

func TestSetCost_MissingService(t *testing.T) {
	ks := keystore.New(tempStorePath(t))
	if err := ks.SetCost("ghost", 100); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestSetCost_NegativeValue(t *testing.T) {
	ks := keystore.New(tempStorePath(t))
	if err := ks.Set("svc", "key"); err != nil {
		t.Fatal(err)
	}
	if err := ks.SetCost("svc", -1); err == nil {
		t.Error("expected error for negative cost")
	}
}

func TestGetCost_MissingService(t *testing.T) {
	ks := keystore.New(tempStorePath(t))
	_, err := ks.GetCost("ghost")
	if err == nil {
		t.Error("expected error for missing service")
	}
}

func TestClearCost_ResetsToZero(t *testing.T) {
	ks := keystore.New(tempStorePath(t))
	if err := ks.Set("svc", "key"); err != nil {
		t.Fatal(err)
	}
	_ = ks.SetCost("svc", 1000)
	if err := ks.ClearCost("svc"); err != nil {
		t.Fatalf("ClearCost: %v", err)
	}
	cost, _ := ks.GetCost("svc")
	if cost != 0 {
		t.Errorf("expected 0 after clear, got %d", cost)
	}
}

func TestServicesByCostAbove_ReturnsMatchingSorted(t *testing.T) {
	ks := keystore.New(tempStorePath(t))
	for _, svc := range []string{"alpha", "beta", "gamma"} {
		_ = ks.Set(svc, "k")
	}
	_ = ks.SetCost("alpha", 500)
	_ = ks.SetCost("beta", 1500)
	_ = ks.SetCost("gamma", 200)

	results := ks.ServicesByCostAbove(400)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0] != "beta" || results[1] != "alpha" {
		t.Errorf("unexpected order: %v", results)
	}
}

func TestSetCost_PersistsAcrossReload(t *testing.T) {
	path := tempStorePath(t)
	ks := keystore.New(path)
	_ = ks.Set("svc", "key")
	_ = ks.SetCost("svc", 750)

	ks2 := keystore.New(path)
	cost, err := ks2.GetCost("svc")
	if err != nil {
		t.Fatalf("GetCost after reload: %v", err)
	}
	if cost != 750 {
		t.Errorf("expected 750 after reload, got %d", cost)
	}
}
