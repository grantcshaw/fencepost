package keystore_test

import (
	"testing"

	"github.com/yourusername/fencepost/internal/keystore"
)

func TestCost_PersistsToDisk(t *testing.T) {
	path := tempStorePath(t)
	ks := keystore.New(path)

	services := []struct {
		name string
		cost int
	}{
		{"aws", 12000},
		{"gcp", 8500},
		{"azure", 3200},
	}
	for _, s := range services {
		if err := ks.Set(s.name, "key"); err != nil {
			t.Fatal(err)
		}
		if err := ks.SetCost(s.name, s.cost); err != nil {
			t.Fatalf("SetCost %s: %v", s.name, err)
		}
	}

	// Reload and verify all costs persisted
	ks2 := keystore.New(path)
	for _, s := range services {
		got, err := ks2.GetCost(s.name)
		if err != nil {
			t.Fatalf("GetCost %s after reload: %v", s.name, err)
		}
		if got != s.cost {
			t.Errorf("%s: expected %d, got %d", s.name, s.cost, got)
		}
	}
}

func TestCost_AboveThreshold_AfterClear(t *testing.T) {
	ks := keystore.New(tempStorePath(t))
	_ = ks.Set("svcA", "k")
	_ = ks.Set("svcB", "k")
	_ = ks.SetCost("svcA", 2000)
	_ = ks.SetCost("svcB", 2000)

	results := ks.ServicesByCostAbove(1000)
	if len(results) != 2 {
		t.Fatalf("expected 2, got %d", len(results))
	}

	_ = ks.ClearCost("svcA")
	results = ks.ServicesByCostAbove(1000)
	if len(results) != 1 {
		t.Fatalf("expected 1 after clear, got %d", len(results))
	}
	if results[0] != "svcB" {
		t.Errorf("expected svcB, got %s", results[0])
	}
}
