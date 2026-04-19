package keystore

import "sort"

// StoreEntry holds all metadata for a single service key.
type StoreEntry struct {
	Key          string            `json:"key"`
	CreatedAt    int64             `json:"created_at"`
	RotatedAt    int64             `json:"rotated_at"`
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
	LastAccessed int64             `json:"last_accessed,omitempty"`
	TTL          int64             `json:"ttl,omitempty"`
	Flags        []string          `json:"flags,omitempty"`
	Credential   string            `json:"credential,omitempty"`
	RotationCount int              `json:"rotation_count,omitempty"`
	Contact      string            `json:"contact,omitempty"`
	URL          string            `json:"url,omitempty"`
	Webhook      string            `json:"webhook,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	Alias        string            `json:"alias,omitempty"`
	Secret       string            `json:"secret,omitempty"`
	Link         string            `json:"link,omitempty"`
	Token        string            `json:"token,omitempty"`
	Quota        int               `json:"quota,omitempty"`
}

type storeData struct {
	Entries map[string]StoreEntry `json:"entries"`
}

func sortStrings(s []string) {
	sort.Strings(s)
}
