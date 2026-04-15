package keystore

import (
	"testing"
	"time"
)

func TestSearch_ByServiceName(t *testing.T) {
	s := newTestStore(t)

	s.Set("github", "key-gh", time.Now())
	s.Set("gitlab", "key-gl", time.Now())
	s.Set("stripe", "key-st", time.Now())

	results := s.Search("git")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Service != "github" || results[1].Service != "gitlab" {
		t.Errorf("unexpected results: %v", results)
	}
}

func TestSearch_ByTag(t *testing.T) {
	s := newTestStore(t)

	s.Set("aws", "key-aws", time.Now())
	s.Set("gcp", "key-gcp", time.Now())
	s.SetTags("aws", []string{"cloud", "prod"})
	s.SetTags("gcp", []string{"cloud", "dev"})

	results := s.Search("cloud")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestSearch_ByNote(t *testing.T) {
	s := newTestStore(t)

	s.Set("sendgrid", "key-sg", time.Now())
	s.Set("mailgun", "key-mg", time.Now())
	s.SetNote("sendgrid", "used for transactional email")

	results := s.Search("transactional")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Service != "sendgrid" {
		t.Errorf("expected sendgrid, got %s", results[0].Service)
	}
}

func TestSearch_NoMatches(t *testing.T) {
	s := newTestStore(t)
	s.Set("github", "key-gh", time.Now())

	results := s.Search("zzznomatch")
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestSearch_CaseInsensitive(t *testing.T) {
	s := newTestStore(t)
	s.Set("GitHub", "key-gh", time.Now())

	results := s.Search("github")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func newTestStore(t *testing.T) *Store {
	t.Helper()
	s, err := New(tempStorePath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}
