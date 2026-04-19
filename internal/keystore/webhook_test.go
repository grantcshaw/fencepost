package keystore

import (
	"testing"
)

func TestSetWebhook_StoresWebhook(t *testing.T) {
	s := newTestStore(t)
	if err := s.Set("svc", "key123"); err != nil {
		t.Fatal(err)
	}
	if err := s.SetWebhook("svc", "https://example.com/hook"); err != nil {
		t.Fatal(err)
	}
	url, err := s.GetWebhook("svc")
	if err != nil {
		t.Fatal(err)
	}
	if url != "https://example.com/hook" {
		t.Errorf("expected webhook URL, got %q", url)
	}
}

func TestSetWebhook_MissingService(t *testing.T) {
	s := newTestStore(t)
	if err := s.SetWebhook("missing", "https://example.com/hook"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestGetWebhook_MissingService(t *testing.T) {
	s := newTestStore(t)
	if _, err := s.GetWebhook("missing"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestClearWebhook_RemovesWebhook(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("svc", "key")
	_ = s.SetWebhook("svc", "https://example.com/hook")
	if err := s.ClearWebhook("svc"); err != nil {
		t.Fatal(err)
	}
	url, _ := s.GetWebhook("svc")
	if url != "" {
		t.Errorf("expected empty webhook, got %q", url)
	}
}

func TestServicesByWebhook_ReturnsMatchingSorted(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("beta", "k1")
	_ = s.Set("alpha", "k2")
	_ = s.Set("gamma", "k3")
	_ = s.SetWebhook("beta", "https://b.example.com")
	_ = s.SetWebhook("alpha", "https://a.example.com")

	result := s.ServicesByWebhook()
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
	if result[0] != "alpha" || result[1] != "beta" {
		t.Errorf("unexpected order: %v", result)
	}
}
