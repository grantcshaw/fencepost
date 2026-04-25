package keystore

import "fmt"

type Visibility string

const (
	VisibilityPublic   Visibility = "public"
	VisibilityPrivate  Visibility = "private"
	VisibilityInternal Visibility = "internal"
)

// isValidVisibility returns true if v is one of the recognized visibility values.
func isValidVisibility(v Visibility) bool {
	return v == VisibilityPublic || v == VisibilityPrivate || v == VisibilityInternal
}

// SetVisibility updates the visibility of the given service entry.
func (ks *KeyStore) SetVisibility(service string, v Visibility) error {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	entry, ok := ks.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}

	if !isValidVisibility(v) {
		return fmt.Errorf("invalid visibility %q: must be public, private, or internal", v)
	}

	entry.Visibility = string(v)
	ks.data.Entries[service] = entry
	return ks.save()
}

// GetVisibility returns the visibility of the given service entry.
// If no visibility has been set, it defaults to VisibilityPrivate.
func (ks *KeyStore) GetVisibility(service string) (Visibility, error) {
	ks.mu.RLock()
	defer ks.mu.RUnlock()

	entry, ok := ks.data.Entries[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}

	if entry.Visibility == "" {
		return VisibilityPrivate, nil
	}
	return Visibility(entry.Visibility), nil
}

// ServicesByVisibility returns a sorted list of service names whose effective
// visibility matches v. Services with no visibility set are treated as private.
func (ks *KeyStore) ServicesByVisibility(v Visibility) ([]string, error) {
	ks.mu.RLock()
	defer ks.mu.RUnlock()

	var results []string
	for name, entry := range ks.data.Entries {
		effective := entry.Visibility
		if effective == "" {
			effective = string(VisibilityPrivate)
		}
		if Visibility(effective) == v {
			results = append(results, name)
		}
	}
	sortResults(results)
	return results, nil
}
