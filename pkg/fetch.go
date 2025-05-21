package seeker

import (
	"context"
	"time"

	gardener_types "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/kyma-project/gardener-syncer/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type List func(context.Context, client.ObjectList, ...client.ListOption) error

type FetchSeeds func() (types.Providers, error)

type FetchSeedsOpts struct {
	Timeout time.Duration
	List
}

func BuildFetchSeedFn(opts FetchSeedsOpts) FetchSeeds {
	return func() (types.Providers, error) {
		ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
		defer logWithDuration(time.Now())
		defer cancel()

		var seeds gardener_types.SeedList
		if err := opts.List(ctx, &seeds); err != nil {
			return nil, err
		}

		return ToProviderRegions(seeds.Items), nil
	}
}
