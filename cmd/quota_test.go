package cmd

import (
	"testing"

	"github.com/fencepost/internal/keystore"
)

func newQuotaTestStore(t *testing.T) *keystore.Store {
	t.Helper()
	s, err := keystore.New(tempStorePath(t))
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestQuotaCmd_SetAndGet(t *testing.T) {
	s := newQuotaTestStore(t)
	if err := s.Set("mysvc", "apikey123"); err != nil {
		t.Fatal(err)
	}
	if err := s.SetQuota("mysvc", 500); err != nil {
		t.Fatal(err)
	}
	q, err := s.GetQuota("mysvc")
	if err != nil {
		t.Fatal(err)
	}
	if q != 500 {
		t.Fatalf("expected 500, got %d", q)
	}
}

func TestQuotaCmd_ListByQuota(t *testing.T) {
	s := newQuotaTestStore(t)
	for _, svc := range []string{"a", "b", "c"} {
		_ = s.Set(svc, "k")
	}
	_ = s.SetQuota("a", 10)
	_ = s.SetQuota("b", 200)
	_ = s.SetQuota("c", 5)

	result := s.ServicesByQuota(10)
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d: %v", len(result), result)
	}
	if result[0] != "a" || result[1] != "b" {
		t.Fatalf("unexpected services: %v", result)
	}
}
