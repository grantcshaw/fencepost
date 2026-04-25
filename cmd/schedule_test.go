package cmd_test

import (
	"bytes"
	"testing"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func newScheduleTestStore(t *testing.T) *keystore.Store {
	t.Helper()
	path := tempStorePath(t)
	s, err := keystore.New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	return s
}

func TestScheduleCmd_SetAndGet(t *testing.T) {
	s := newScheduleTestStore(t)
	if err := s.Set("myapi", "secret-key"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if err := s.SetSchedule("myapi", "48h"); err != nil {
		t.Fatalf("SetSchedule failed: %v", err)
	}

	got, err := s.GetSchedule("myapi")
	if err != nil {
		t.Fatalf("GetSchedule failed: %v", err)
	}
	if got != "48h" {
		t.Errorf("expected 48h, got %q", got)
	}
}

func TestScheduleCmd_ClearSchedule(t *testing.T) {
	s := newScheduleTestStore(t)
	_ = s.Set("myapi", "secret-key")
	_ = s.SetSchedule("myapi", "24h")

	if err := s.ClearSchedule("myapi"); err != nil {
		t.Fatalf("ClearSchedule failed: %v", err)
	}
	got, _ := s.GetSchedule("myapi")
	if got != "" {
		t.Errorf("expected empty schedule after clear, got %q", got)
	}
}

func TestScheduleCmd_ListBySchedule(t *testing.T) {
	s := newScheduleTestStore(t)
	for _, svc := range []string{"alpha", "beta", "gamma"} {
		_ = s.Set(svc, "k")
	}
	_ = s.SetSchedule("alpha", "24h")
	_ = s.SetSchedule("gamma", "168h")

	var buf bytes.Buffer
	results := s.ServicesBySchedule()
	for _, r := range results {
		buf.WriteString(r + "\n")
	}

	output := buf.String()
	if output == "" {
		t.Error("expected non-empty output")
	}
	if len(results) != 2 {
		t.Errorf("expected 2 scheduled services, got %d", len(results))
	}
}
