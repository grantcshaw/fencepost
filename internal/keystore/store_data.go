package keystore

import "time"

// storeData is the serialized form of the store on disk.
// Extending here keeps keystore.go clean and centralises schema changes.
type storeDataWatchExtension struct {
	Watched map[string]time.Time `json:"watched,omitempty"`
}
