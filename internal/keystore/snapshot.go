package keystore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot represents a point-in-time backup of the keystore.
type Snapshot struct {
	CreatedAt time.Time         `json:"created_at"`
	Entries   map[string]Record `json:"entries"`
}

// WriteSnapshot saves a snapshot of the current keystore to the given directory.
// The file is named snapshot-<unix-timestamp>.json.
func (s *Store) WriteSnapshot(dir string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("create snapshot dir: %w", err)
	}

	snap := Snapshot{
		CreatedAt: time.Now().UTC(),
		Entries:   make(map[string]Record, len(s.data.Records)),
	}
	for k, v := range s.data.Records {
		snap.Entries[k] = v
	}

	filename := fmt.Sprintf("snapshot-%d.json", snap.CreatedAt.Unix())
	path := filepath.Join(dir, filename)

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0600)
	if err != nil {
		return "", fmt.Errorf("open snapshot file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		return "", fmt.Errorf("encode snapshot: %w", err)
	}

	return path, nil
}

// LoadSnapshot reads a snapshot file and returns a Snapshot.
func LoadSnapshot(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open snapshot: %w", err)
	}
	defer f.Close()

	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return nil, fmt.Errorf("decode snapshot: %w", err)
	}
	return &snap, nil
}
