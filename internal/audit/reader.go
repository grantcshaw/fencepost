package audit

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// ReadAll reads all audit log entries from the given file path.
func ReadAll(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Entry{}, nil
		}
		return nil, fmt.Errorf("audit: open log file: %w", err)
	}
	defer f.Close()

	var entries []Entry
	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var entry Entry
		if err := json.Unmarshal(line, &entry); err != nil {
			return nil, fmt.Errorf("audit: parse line %d: %w", lineNum, err)
		}
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("audit: scan log file: %w", err)
	}
	return entries, nil
}

// FilterByService returns entries that match the given service name.
func FilterByService(entries []Entry, service string) []Entry {
	var result []Entry
	for _, e := range entries {
		if e.Service == service {
			result = append(result, e)
		}
	}
	return result
}

// FilterByEvent returns entries that match the given event type.
func FilterByEvent(entries []Entry, event EventType) []Entry {
	var result []Entry
	for _, e := range entries {
		if e.Event == event {
			result = append(result, e)
		}
	}
	return result
}
