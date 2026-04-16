package keystore

import "fmt"

// MergeResult describes what happened to a single service during a merge.
type MergeResult struct {
	Service   string
	Action    string // "added", "skipped", "overwritten"
}

// MergeOptions controls merge behaviour.
type MergeOptions struct {
	Overwrite bool
}

// Merge imports all entries from src into ks.
// Returns a slice of MergeResult describing each service processed.
func (ks *KeyStore) Merge(src *KeyStore, opts MergeOptions) ([]MergeResult, error) {
	names, err := src.List()
	if err != nil {
		return nil, fmt.Errorf("merge: list source: %w", err)
	}

	var results []MergeResult
	for _, name := range names {
		entry, err := src.Get(name)
		if err != nil {
			return nil, fmt.Errorf("merge: get %q: %w", name, err)
		}

		_, missing := ks.Get(name)
		exists := missing == nil

		if exists && !opts.Overwrite {
			results = append(results, MergeResult{Service: name, Action: "skipped"})
			continue
		}

		if err := ks.Set(name, entry.Key); err != nil {
			return nil, fmt.Errorf("merge: set %q: %w", name, err)
		}

		// Carry over tags and note.
		tags, _ := src.GetTags(name)
		if len(tags) > 0 {
			_ = ks.SetTags(name, tags)
		}
		note, _ := src.GetNote(name)
		if note != "" {
			_ = ks.SetNote(name, note)
		}

		action := "added"
		if exists {
			action = "overwritten"
		}
		results = append(results, MergeResult{Service: name, Action: action})
	}
	return results, nil
}
