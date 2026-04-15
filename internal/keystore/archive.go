package keystore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ArchiveEntry holds a removed service entry along with the time it was archived.
type ArchiveEntry struct {
	ArchivedAt time.Time  `json:"archived_at"`
	Entry      Entry      `json:"entry"`
}

// Archive moves a service from the active store to an archive file and removes
// it from the live keystore. Returns an error if the service does not exist.
func (s *Store) Archive(service, archivePath string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}

	existing, err := loadArchiveFile(archivePath)
	if err != nil {
		return fmt.Errorf("loading archive: %w", err)
	}

	existing = append(existing, ArchiveEntry{
		ArchivedAt: time.Now().UTC(),
		Entry:      entry,
	})

	if err := writeArchiveFile(archivePath, existing); err != nil {
		return fmt.Errorf("writing archive: %w", err)
	}

	delete(s.data, service)
	return s.persist()
}

// LoadArchive reads all archived entries from the given file path.
func LoadArchive(archivePath string) ([]ArchiveEntry, error) {
	return loadArchiveFile(archivePath)
}

func loadArchiveFile(path string) ([]ArchiveEntry, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []ArchiveEntry{}, nil
	}
	if err != nil {
		return nil, err
	}
	var entries []ArchiveEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func writeArchiveFile(path string, entries []ArchiveEntry) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}
