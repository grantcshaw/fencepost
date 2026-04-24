package cmd_test

import (
	"testing"

	"github.com/yourusername/fencepost/internal/keystore"
)

func newCostTestStore(t *testing.T) (string, *keystore.Store) {
	t.Helper()
	path := t.TempDir() + "/store.json"
	ks := keystore.New(path)
	return path, ks
}

func TestCostCmd_SetAndGet(t *testing.T) {
	path, ks := newCostTestStore(t)
	if err := ks.Set("stripe", "sk_live_abc"); err != nil {
		t.Fatal(err)
	}
	if err := ks.SetCost("stripe", 2500); err != nil {
		t.Fatalf("SetCost: %v", err)
	}

	ks2 := keystore.New(path)
	cost, err := ks2.GetCost("stripe")
	if err != nil {
		t.Fatalf("GetCost: %v", err)
	}
	if cost != 2500 {
		t.Errorf("expected 2500, got %d", cost)
	}
}

func TestCostCmd_ListAbove(t *testing.T) {
	_, ks := newCostTestStore(t)
	services := map[string]int{
		"cheap":    50,
		"moderate": 999,
		"expensive": 5000,
	}
	for svc, cost := range services {
		_ = ks.Set(svc, "key")
		_ = ks.SetCost(svc, cost)
	}

	results := ks.ServicesByCostAbove(500)
	if len(results) != 2 {
		t.Fatalf("expected 2 results above 500 cents, got %d: %v", len(results), results)
	}
	if results[0] != "expensive" {
		t.Errorf("expected expensive first (highest cost), got %s", results[0])
	}
}
