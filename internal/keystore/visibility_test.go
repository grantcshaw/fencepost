package keystore

import (
	"testing"
)

func TestSetVisibility_StoresVisibility(t *testing.T) {
	ks := New(tempStorePath(t))
	ks.Set("svc", "key123")

	if err := ks.SetVisibility("svc", VisibilityPublic); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	v, err := ks.GetVisibility("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != VisibilityPublic {
		t.Errorf("expected public, got %q", v)
	}
}

func TestSetVisibility_MissingService(t *testing.T) {
	ks := New(tempStorePath(t))
	err := ks.SetVisibility("ghost", VisibilityPublic)
	if err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestSetVisibility_InvalidValue(t *testing.T) {
	ks := New(tempStorePath(t))
	ks.Set("svc", "key123")
	err := ks.SetVisibility("svc", Visibility("secret"))
	if err == nil {
		t.Fatal("expected error for invalid visibility")
	}
}

func TestGetVisibility_DefaultsToPrivate(t *testing.T) {
	ks := New(tempStorePath(t))
	ks.Set("svc", "key123")

	v, err := ks.GetVisibility("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != VisibilityPrivate {
		t.Errorf("expected private default, got %q", v)
	}
}

func TestGetVisibility_MissingService(t *testing.T) {
	ks := New(tempStorePath(t))
	_, err := ks.GetVisibility("ghost")
	if err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestServicesByVisibility_ReturnsMatchingSorted(t *testing.T) {
	ks := New(tempStorePath(t))
	for _, s := range []string{"alpha", "beta", "gamma"} {
		ks.Set(s, "k")
	}
	ks.SetVisibility("alpha", VisibilityInternal)
	ks.SetVisibility("gamma", VisibilityInternal)

	results, err := ks.ServicesByVisibility(VisibilityInternal)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 || results[0] != "alpha" || results[1] != "gamma" {
		t.Errorf("unexpected results: %v", results)
	}
}
