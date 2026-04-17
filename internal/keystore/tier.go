package keystore

import "fmt"

type Tier string

const (
	TierFree       Tier = "free"
	TierBasic      Tier = "basic"
	TierPro        Tier = "pro"
	TierEnterprise Tier = "enterprise"
)

var validTiers = map[Tier]bool{
	TierFree: true, TierBasic: true, TierPro: true, TierEnterprise: true,
}

func (ks *KeyStore) SetTier(service string, tier Tier) error {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	entry, ok := ks.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	if !validTiers[tier] {
		return fmt.Errorf("invalid tier %q: must be one of free, basic, pro, enterprise", tier)
	}
	entry.Tier = string(tier)
	ks.data.Entries[service] = entry
	return ks.save()
}

func (ks *KeyStore) GetTier(service string) (Tier, error) {
	ks.mu.RLock()
	defer ks.mu.RUnlock()
	entry, ok := ks.data.Entries[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}
	return Tier(entry.Tier), nil
}

func (ks *KeyStore) ClearTier(service string) error {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	entry, ok := ks.data.Entries[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	entry.Tier = ""
	ks.data.Entries[service] = entry
	return ks.save()
}

func (ks *KeyStore) ServicesByTier(tier Tier) []string {
	ks.mu.RLock()
	defer ks.mu.RUnlock()
	var results []string
	for name, entry := range ks.data.Entries {
		if Tier(entry.Tier) == tier {
			results = append(results, name)
		}
	}
	sortResults(results)
	return results
}
