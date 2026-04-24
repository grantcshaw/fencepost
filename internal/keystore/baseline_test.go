package keystore_test

import (
	"testing"
	"time"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func TestSetBaseline_StoresBaseline(t *testing.T) {
	s := New(tempStorePath(t))
	_ = s.Set("svc", "key-abc")

	if err := s.SetBaseline("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	baseline, ts, err := s.GetBaseline("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if baseline != "key-abc" {
		t.Errorf("expected baseline %q, got %q", "key-abc", baseline)
	}
	if ts.IsZero() {
		t.Error("expected non-zero BaselineAt timestamp")
	}
}

func TestSetBaseline_MissingService(t *testing.T) {
	s := New(tempStorePath(t))
	if err := s.SetBaseline("ghost"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestGetBaseline_MissingService(t *testing.T) {
	s := New(tempStorePath(t))
	_, _, err := s.GetBaseline("ghost")
	if err == nil {
		t.Error("expected error for missing service")
	}
}

func TestClearBaseline_RemovesBaseline(t *testing.T) {
	s := New(tempStorePath(t))
	_ = s.Set("svc", "key-abc")
	_ = s.SetBaseline("svc")

	if err := s.ClearBaseline("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	baseline, ts, _ := s.GetBaseline("svc")
	if baseline != "" {
		t.Errorf("expected empty baseline, got %q", baseline)
	}
	if !ts.IsZero() {
		t.Error("expected zero BaselineAt after clear")
	}
}

func TestBaselineChanged_DetectsDrift(t *testing.T) {
	s := New(tempStorePath(t))
	_ = s.Set("svc", "key-original")
	_ = s.SetBaseline("svc")

	changed, err := s.BaselineChanged("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if changed {
		t.Error("expected no change immediately after baseline set")
	}

	_ = s.Set("svc", "key-rotated")
	changed, err = s.BaselineChanged("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !changed {
		t.Error("expected change detected after key rotation")
	}
}

func TestServicesWithBaseline_ReturnsSorted(t *testing.T) {
	s := New(tempStorePath(t))
	for _, svc := range []string{"gamma", "alpha", "beta"} {
		_ = s.Set(svc, "k")
		_ = s.SetBaseline(svc)
	}
	_ = s.Set("delta", "k") // no baseline

	names := s.ServicesWithBaseline()
	expected := []string{"alpha", "beta", "gamma"}
	if len(names) != len(expected) {
		t.Fatalf("expected %d services, got %d", len(expected), len(names))
	}
	for i, name := range names {
		if name != expected[i] {
			t.Errorf("index %d: expected %q, got %q", i, expected[i], name)
		}
	}
	_ = time.Now() // suppress unused import
}
