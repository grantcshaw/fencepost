package keystore

import "sort"

// StoreData holds all persisted keystore state.
type StoreData struct {
	Entries map[string]Entry `json:"entries"`
}

// Entry represents a single stored API key with all its metadata.
type Entry struct {
	Key          string            `json:"key"`
	CreatedAt    int64             `json:"created_at"`
	RotatedAt    int64             `json:"rotated_at,omitempty"`
	Tags         []string          `json:"tags,omitempty"`
	Note         string            `json:"note,omitempty"`
	Label        string            `json:"label,omitempty"`
	Group        string            `json:"group,omitempty"`
	Pinned       bool              `json:"pinned,omitempty"`
	Watched      bool              `json:"watched,omitempty"`
	Priority     string            `json:"priority,omitempty"`
	Comment      string            `json:"comment,omitempty"`
	Environment  string            `json:"environment,omitempty"`
	Visibility   string            `json:"visibility,omitempty"`
	Owner        string            `json:"owner,omitempty"`
	Category     string            `json:"category,omitempty"`
	Version      string            `json:"version,omitempty"`
	Region       string            `json:"region,omitempty"`
	Tier         string            `json:"tier,omitempty"`
	Lifecycle    string            `json:"lifecycle,omitempty"`
	Dependencies []string          `json:"dependencies,omitempty"`
	Scope        string            `json:"scope,omitempty"`
	Rating       string            `json:"rating,omitempty"`
	Source       string            `json:"source,omitempty"`
	Protocol     string            `json:"protocol,omitempty"`
	Format       string            `json:"format,omitempty"`
	LastAccessed int64             `json:"last_accessed,omitempty"`
	TTL          int64             `json:"ttl,omitempty"`
	Flags        []string          `json:"flags,omitempty"`
	Meta         map[string]string `json:"meta,omitempty"`
}

func sortStrings(s []string) {
	sort.Strings(s)
}
