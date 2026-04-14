package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// EventType represents the type of audit event.
type EventType string

const (
	EventKeyCreated EventType = "key_created"
	EventKeyRotated EventType = "key_rotated"
	EventKeyRevoked EventType = "key_revoked"
	EventKeyListed  EventType = "key_listed"
)

// Entry represents a single audit log entry.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     EventType `json:"event"`
	Service   string    `json:"service"`
	KeyID     string    `json:"key_id,omitempty"`
	Message   string    `json:"message"`
}

// Logger writes audit entries to a file in JSON Lines format.
type Logger struct {
	path string
}

// New creates a new Logger that writes to the given file path.
func New(path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("audit: open log file: %w", err)
	}
	f.Close()
	return &Logger{path: path}, nil
}

// Log writes an audit entry to the log file.
func (l *Logger) Log(event EventType, service, keyID, message string) error {
	entry := Entry{
		Timestamp: time.Now().UTC(),
		Event:     event,
		Service:   service,
		KeyID:     keyID,
		Message:   message,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}

	f, err := os.OpenFile(l.path, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("audit: open log file: %w", err)
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "%s\n", data)
	if err != nil {
		return fmt.Errorf("audit: write entry: %w", err)
	}
	return nil
}
