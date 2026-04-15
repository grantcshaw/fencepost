package keystore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// BackupMeta holds metadata about a backup file.
type BackupMeta struct {
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
	Services  int       `json:"services"`
}

// WriteBackup writes a timestamped JSON backup of the store to the given directory.
// Returns metadata about the written backup.
func (s *Store) WriteBackup(dir string) (BackupMeta, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if err := os.MkdirAll(dir, 0700); err != nil {
		return BackupMeta{}, fmt.Errorf("create backup dir: %w", err)
	}

	timestamp := time.Now().UTC().Format("20060102T150405Z")
	filename := fmt.Sprintf("fencepost-backup-%s.json", timestamp)
	path := filepath.Join(dir, filename)

	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return BackupMeta{}, fmt.Errorf("marshal backup: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return BackupMeta{}, fmt.Errorf("write backup file: %w", err)
	}

	return BackupMeta{
		Path:      path,
		CreatedAt: time.Now().UTC(),
		Services:  len(s.data),
	}, nil
}

// RestoreBackup loads store data from a backup file, replacing all current entries.
// The store is persisted to its configured path after restore.
func (s *Store) RestoreBackup(backupPath string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	raw, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("read backup file: %w", err)
	}

	var restored map[string]Entry
	if err := json.Unmarshal(raw, &restored); err != nil {
		return fmt.Errorf("parse backup file: %w", err)
	}

	s.data = restored
	return s.saveLocked()
}
