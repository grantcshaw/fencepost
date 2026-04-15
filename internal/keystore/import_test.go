package keystore_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

func writeImportFile(t *testing.T, records []keystore.ImportRecord) string {
	t.Helper()
	data, err := json.Marshal(records)
	if err != nil {
		t.Fatalf("marshal import file: %v", err)
	}
	p := filepath.Join(t.TempDir(), "import.json")
	if err := os.WriteFile(p, data, 0600); err != nil {
		t.Fatalf("write import file: %v", err)
	}
	return p
}

func TestImportFromFile_ImportsAllRecords(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	records := []keystore.ImportRecord{
		{Service: "alpha", Key: "key-alpha"},
		{Service: "beta", Key: "key-beta", Tags: []string{"prod"}, Note: "main key"},
	}
	path := writeImportFile(t, records)

	res, err := s.ImportFromFile(path, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Imported) != 2 {
		t.Errorf("expected 2 imported, got %d", len(res.Imported))
	}
	if len(res.Skipped) != 0 {
		t.Errorf("expected 0 skipped, got %d", len(res.Skipped))
	}

	v, ok := s.Get("beta")
	if !ok || v != "key-beta" {
		t.Errorf("expected key-beta, got %q (ok=%v)", v, ok)
	}
}

func TestImportFromFile_SkipsExistingWithoutOverwrite(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	_ = s.Set("alpha", "original")

	records := []keystore.ImportRecord{
		{Service: "alpha", Key: "new-value"},
	}
	path := writeImportFile(t, records)

	res, err := s.ImportFromFile(path, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "alpha" {
		t.Errorf("expected alpha in skipped, got %v", res.Skipped)
	}

	v, _ := s.Get("alpha")
	if v != "original" {
		t.Errorf("expected original key to be preserved, got %q", v)
	}
}

func TestImportFromFile_OverwritesExisting(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	_ = s.Set("alpha", "original")

	records := []keystore.ImportRecord{
		{Service: "alpha", Key: "updated"},
	}
	path := writeImportFile(t, records)

	res, err := s.ImportFromFile(path, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Imported) != 1 {
		t.Errorf("expected 1 imported, got %d", len(res.Imported))
	}

	v, _ := s.Get("alpha")
	if v != "updated" {
		t.Errorf("expected updated key, got %q", v)
	}
}

func TestImportFromFile_MissingFile(t *testing.T) {
	s := keystore.New(tempStorePath(t))
	_, err := s.ImportFromFile("/nonexistent/path.json", false)
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
