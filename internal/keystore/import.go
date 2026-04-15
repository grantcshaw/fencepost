package keystore

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// ImportRecord represents a single service entry in an import file.
type ImportRecord struct {
	Service   string    `json:"service"`
	Key       string    `json:"key"`
	Tags      []string  `json:"tags,omitempty"`
	Note      string    `json:"note,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// ImportResult summarises the outcome of a bulk import operation.
type ImportResult struct {
	Imported  []string
	Skipped   []string
	Overwrite bool
}

// ImportFromFile reads a JSON file containing an array of ImportRecord values
// and stores each one in the keystore. When overwrite is false, existing keys
// are left unchanged and their service names are collected in Skipped.
func (s *Store) ImportFromFile(path string, overwrite bool) (*ImportResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("import: read file: %w", err)
	}

	var records []ImportRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("import: parse JSON: %w", err)
	}

	result := &ImportResult{Overwrite: overwrite}

	for _, r := range records {
		if r.Service == "" || r.Key == "" {
			continue
		}

		_, exists := s.Get(r.Service)
		if exists && !overwrite {
			result.Skipped = append(result.Skipped, r.Service)
			continue
		}

		if err := s.Set(r.Service, r.Key); err != nil {
			return nil, fmt.Errorf("import: set %q: %w", r.Service, err)
		}

		if len(r.Tags) > 0 {
			if _, err := s.SetTags(r.Service, r.Tags); err != nil {
				return nil, fmt.Errorf("import: set tags for %q: %w", r.Service, err)
			}
		}

		if r.Note != "" {
			if err := s.SetNote(r.Service, r.Note); err != nil {
				return nil, fmt.Errorf("import: set note for %q: %w", r.Service, err)
			}
		}

		result.Imported = append(result.Imported, r.Service)
	}

	return result, nil
}
