package keystore

// StoreEntry holds all metadata for a single service key.
type StoreEntry struct {
	Key          string   `json:"key"`
	CreatedAt    int64    `json:"created_at"`
	RotatedAt    int64    `json:"rotated_at"`
	Tags         []string `json:"tags,omitempty"`
	Note         string   `json:"note,omitempty"`
	Label        string   `json:"label,omitempty"`
	Group        string   `json:"group,omitempty"`
	Pinned       bool     `json:"pinned,omitempty"`
	Watched      bool     `json:"watched,omitempty"`
	Priority     string   `json:"priority,omitempty"`
	Comment      string   `json:"comment,omitempty"`
	Environment  string   `json:"environment,omitempty"`
	Visibility   string   `json:"visibility,omitempty"`
	Owner        string   `json:"owner,omitempty"`
	Category     string   `json:"category,omitempty"`
	Version      string   `json:"version,omitempty"`
	Region       string   `json:"region,omitempty"`
	Tier         string   `json:"tier,omitempty"`
	Lifecycle    string   `json:"lifecycle,omitempty"`
	Dependencies []string `json:"dependencies,omitempty"`
	Scope        string   `json:"scope,omitempty"`
	Rating       string   `json:"rating,omitempty"`
	Source       string   `json:"source,omitempty"`
}

// storeData is the top-level structure persisted to disk.
type storeData struct {
	Entries map[string]StoreEntry `json:"entries"`
}
