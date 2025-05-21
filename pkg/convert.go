package seeker

import (
	"strings"

	"sigs.k8s.io/yaml"

	gardener_types "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	v1beta1helper "github.com/gardener/gardener/pkg/apis/core/v1beta1/helper"
	"github.com/kyma-project/gardener-syncer/pkg/types"
)

func verifySeedReadiness(seed *gardener_types.Seed) bool {
	if seed.Status.LastOperation == nil {
		return false
	}

	if cond := v1beta1helper.GetCondition(seed.Status.Conditions, gardener_types.SeedGardenletReady); cond == nil || cond.Status != gardener_types.ConditionTrue {
		return false
	}

	if seed.Spec.Backup != nil {
		if cond := v1beta1helper.GetCondition(seed.Status.Conditions, gardener_types.SeedBackupBucketsReady); cond == nil || cond.Status != gardener_types.ConditionTrue {
			return false
		}
	}

	return true
}

func seedCanBeUsed(seed *gardener_types.Seed) bool {
	return seed.DeletionTimestamp == nil && seed.Spec.Settings.Scheduling.Visible && verifySeedReadiness(seed)
}

func ToProviderRegions(seeds []gardener_types.Seed) types.Providers {
	result := types.Providers{}
	for _, seed := range seeds {
		if seedCanBeUsed(&seed) {
			result.Add(
				seed.Spec.Provider.Type,
				seed.Spec.Provider.Region,
			)
		}
	}

	return result
}

func ToConfigMap(providerRegions types.Providers) (map[string]string, error) {
	result := map[string]string{}
	for k, v := range providerRegions {
		data, err := yaml.Marshal(v)
		if err != nil {
			return nil, err
		}
		result[k] = strings.TrimRight(string(data), "\n")
	}
	return result, nil
}

type Convert[T any, V any] func(T) (V, error)
