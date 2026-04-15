package keystore

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// ExportEntry represents a single key entry for export.
type ExportEntry struct {
	Service   string   `json:"service"`
	Key       string   `json:"key"`
	Tags      []string `json:"tags,omitempty"`
	Note      string   `json:"note,omitempty"`
	CreatedAt int64    `json:"created_at"`
	RotatedAt int64    `json:"rotated_at,omitempty"`
}

// ExportAll serialises all keys in the store to JSON and writes to path.
func (s *Store) ExportAll(path string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries := make([]ExportEntry, 0, len(s.data.Keys))
	for name, record := range s.data.Keys {
		e := ExportEntry{
			Service:   name,
			Key:       record.Key,
			Tags:      record.Tags,
			Note:      record.Note,
			CreatedAt: record.CreatedAt,
			RotatedAt: record.RotatedAt,
		}
		entries = append(entries, e)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Service < entries[j].Service
	})

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("export: create file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(entries); err != nil {
		return fmt.Errorf("export: encode: %w", err)
	}
	return nil
}

// ExportServices serialises only the named services to JSON and writes to path.
func (s *Store) ExportServices(path string, services []string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries := make([]ExportEntry, 0, len(services))
	for _, name := range services {
		record, ok := s.data.Keys[name]
		if !ok {
			return fmt.Errorf("export: service %q not found", name)
		}
		entries = append(entries, ExportEntry{
			Service:   name,
			Key:       record.Key,
			Tags:      record.Tags,
			Note:      record.Note,
			CreatedAt: record.CreatedAt,
			RotatedAt: record.RotatedAt,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Service < entries[j].Service
	})

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("export: create file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(entries); err != nil {
		return fmt.Errorf("export: encode: %w", err)
	}
	return nil
}
