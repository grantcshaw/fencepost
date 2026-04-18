package keystore

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

// ErrKeyNotFound is returned when a key for a given service is not found.
var ErrKeyNotFound = errors.New("key not found")

// Entry represents a stored API key entry for a service.
type Entry struct {
	Service   string    `json:"service"`
	Key       string    `json:"key"`
	CreatedAt time.Time `json:"created_at"`
	RotatedAt time.Time `json:"rotated_at,omitempty"`
}

// Store manages API key entries persisted to a JSON file.
type Store struct {
	mu      sync.RWMutex
	path    string
	entries map[string]Entry
}

// New loads (or creates) a key store at the given file path.
func New(path string) (*Store, error) {
	s := &Store{
		path:    path,
		entries: make(map[string]Entry),
	}
	if err := s.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return s, nil
}

// Set adds or updates the API key for a service.
func (s *Store) Set(service, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	existing, exists := s.entries[service]
	entry := Entry{
		Service:   service,
		Key:       key,
		CreatedAt: now,
	}
	if exists {
		entry.CreatedAt = existing.CreatedAt
		entry.RotatedAt = now
	}
	s.entries[service] = entry
	return s.save()
}

// Get retrieves the API key entry for a service.
func (s *Store) Get(service string) (Entry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.entries[service]
	if !ok {
		return Entry{}, ErrKeyNotFound
	}
	return entry, nil
}

// Delete removes the API key entry for a service.
func (s *Store) Delete(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.entries[service]; !ok {
		return ErrKeyNotFound
	}
	delete(s.entries, service)
	return s.save()
}

// List returns all stored entries.
func (s *Store) List() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		list = append(list, e)
	}
	return list
}

// Has reports whether an entry exists for the given service.
func (s *Store) Has(service string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.entries[service]
	return ok
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.entries)
}

func (s *Store) save() error {
	data, err := json.MarshalIndent(s.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0600)
}
