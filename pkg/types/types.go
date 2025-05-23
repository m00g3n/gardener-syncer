package types

import "slices"

type ProviderInfo struct {
	SeedRegions []string `json:"seedRegions"`
}

type Providers map[string]ProviderInfo

func (s *Providers) Add(provider, regionName string) {
	providerInfo := (*s)[provider]
	if slices.Contains(providerInfo.SeedRegions, regionName) {
		return
	}

	providerInfo.SeedRegions = append(providerInfo.SeedRegions, regionName)
	(*s)[provider] = providerInfo
}
