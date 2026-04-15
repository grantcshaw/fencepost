package keystore

import "strings"

// SearchResult holds a matched service name and its associated key entry.
type SearchResult struct {
	Service string
	Entry   KeyEntry
}

// Search returns all services whose name or tags contain the given query
// (case-insensitive substring match).
func (s *Store) Search(query string) []SearchResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	q := strings.ToLower(query)
	var results []SearchResult

	for name, entry := range s.data.Keys {
		if matchesSearch(name, entry, q) {
			results = append(results, SearchResult{
				Service: name,
				Entry:   entry,
			})
		}
	}

	sortResults(results)
	return results
}

func matchesSearch(name string, entry KeyEntry, q string) bool {
	if strings.Contains(strings.ToLower(name), q) {
		return true
	}
	for _, tag := range entry.Tags {
		if strings.Contains(strings.ToLower(tag), q) {
			return true
		}
	}
	if strings.Contains(strings.ToLower(entry.Note), q) {
		return true
	}
	return false
}

func sortResults(results []SearchResult) {
	for i := 1; i < len(results); i++ {
		for j := i; j > 0 && results[j].Service < results[j-1].Service; j-- {
			results[j], results[j-1] = results[j-1], results[j]
		}
	}
}
