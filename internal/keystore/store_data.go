package keystore

// StoreData holds all persisted keystore entries.
type StoreData struct {
	Entries map[string]Entry `json:"entries"`
}

// Entry represents a single service key record.
type Entry struct {
	Key         string   `json:"key"`
	CreatedAt   int64    `json:"created_at"`
	RotatedAt   int64    `json:"rotated_at"`
	Tags        []string `json:"tags,omitempty"`
	Note        string   `json:"note,omitempty"`
	Label       string   `json:"label,omitempty"`
	Group       string   `json:"group,omitempty"`
	Pinned      bool     `json:"pinned,omitempty"`
	Watched     bool     `json:"watched,omitempty"`
	Priority    string   `json:"priority,omitempty"`
	Comment     string   `json:"comment,omitempty"`
	Environment string   `json:"environment,omitempty"`
	Visibility  string   `json:"visibility,omitempty"`
	Owner       string   `json:"owner,omitempty"`
	Category    string   `json:"category,omitempty"`
	Version     string   `json:"version,omitempty"`
	Region      string   `json:"region,omitempty"`
}
