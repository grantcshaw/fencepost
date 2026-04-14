package cmd

import (
	"crypto/rand"
	"encoding/hex"
)

const keyBytes = 32

// generateKey returns a cryptographically random hex-encoded API key.
func generateKey() string {
	b := make([]byte, keyBytes)
	if _, err := rand.Read(b); err != nil {
		// rand.Read only fails on catastrophic OS-level errors; panic is appropriate.
		panic("fencepost: failed to read random bytes: " + err.Error())
	}
	return hex.EncodeToString(b)
}
