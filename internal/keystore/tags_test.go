package keystore

import (
	"testing"
)

func TestSetTags_StoresAndReturnsSorted(t *testing.T) {
	path := tempStorePath(t)
	s, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	_ = s.Set("svc-a", "key-a")
	if err := s.SetTags("svc-a", []string{"prod", "api", "critical"}); err != nil {
		t.Fatalf("SetTags: %v", err)
	}

	tags, err := s.GetTags("svc-a")
	if err != nil {
		t.Fatalf("GetTags: %v", err)
	}
	if len(tags) != 3 || tags[0] != "api" || tags[1] != "critical" || tags[2] != "prod" {
		t.Errorf("unexpected tags: %v", tags)
	}
}

func TestSetTags_MissingService(t *testing.T) {
	path := tempStorePath(t)
	s, _ := New(path)

	err := s.SetTags("nonexistent", []string{"x"})
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestGetTags_MissingService(t *testing.T) {
	path := tempStorePath(t)
	s, _ := New(path)

	_, err := s.GetTags("ghost")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestFilterByTag_ReturnsMatchingSorted(t *testing.T) {
	path := tempStorePath(t)
	s, _ := New(path)

	_ = s.Set("svc-b", "key-b")
	_ = s.Set("svc-a", "key-a")
	_ = s.Set("svc-c", "key-c")
	_ = s.SetTags("svc-a", []string{"prod"})
	_ = s.SetTags("svc-b", []string{"prod", "staging"})
	_ = s.SetTags("svc-c", []string{"staging"})

	matches := s.FilterByTag("prod")
	if len(matches) != 2 || matches[0] != "svc-a" || matches[1] != "svc-b" {
		t.Errorf("unexpected matches: %v", matches)
	}
}

func TestFilterByTag_NoMatches(t *testing.T) {
	path := tempStorePath(t)
	s, _ := New(path)

	_ = s.Set("svc-a", "key-a")
	_ = s.SetTags("svc-a", []string{"dev"})

	matches := s.FilterByTag("prod")
	if len(matches) != 0 {
		t.Errorf("expected no matches, got %v", matches)
	}
}

func TestSetTags_PersistsAcrossReload(t *testing.T) {
	path := tempStorePath(t)
	s, _ := New(path)

	_ = s.Set("svc-a", "key-a")
	_ = s.SetTags("svc-a", []string{"infra"})

	s2, err := New(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	tags, err := s2.GetTags("svc-a")
	if err != nil || len(tags) != 1 || tags[0] != "infra" {
		t.Errorf("tags not persisted: %v %v", tags, err)
	}
}
