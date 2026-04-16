package keystore

import "fmt"

type Visibility string

const (
	VisibilityPublic   Visibility = "public"
	VisibilityPrivate  Visibility = "private"
	VisibilityInternal Visibility = "internal"
)

func (ks *KeyStore) SetVisibility(service string, v Visibility) error {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	entry, ok := ks.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}

	if v != VisibilityPublic && v != VisibilityPrivate && v != VisibilityInternal {
		return fmt.Errorf("invalid visibility %q: must be public, private, or internal", v)
	}

	entry.Visibility = string(v)
	ks.data.Entries[service] = entry
	return ks.save()
}

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
