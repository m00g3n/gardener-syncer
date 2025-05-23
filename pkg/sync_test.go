package seeker_test

import (
	"fmt"
	"testing"

	seeker "github.com/kyma-project/gardener-syncer/pkg"
	"github.com/kyma-project/gardener-syncer/pkg/types"
	"github.com/stretchr/testify/require"
)

var (
	errStoreFailedTest       = fmt.Errorf("store test failed")
	errFetchSeedsFailedTest = fmt.Errorf("fetch seeds test fail")
)

func TestBuildSyncFn(t *testing.T) {
	testCases := []struct {
		name        string
		store       seeker.Store
		fetch       seeker.FetchSeeds
		expectedErr error
	}{
		{
			name:        "fetch error",
			fetch:       buildFetchSeedsWithError(errFetchSeedsFailedTest),
			expectedErr: errFetchSeedsFailedTest,
		},
		{
			name:        "store error",
			fetch:       buildFetch(types.Providers{}),
			store:       buildStoreWithError(errStoreFailedTest),
			expectedErr: errSoreFailedTest,
		},
		{
			name:  "OK",
			fetch: buildFetch(types.Providers{}),
			store: buildStore(),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// GIVEN
			sync := seeker.BuildSyncFn(testCase.store, testCase.fetch)

			// WHEN
			err := sync()

			// THEN
			if testCase.expectedErr == nil {
				require.NoError(t, err)
			}

			// THEN
			if testCase.expectedErr != nil {
				require.EqualError(t, err, testCase.expectedErr.Error())
			}
		})
	}
}

func buildFetchSeedsWithError(err error) seeker.FetchSeeds {
	return func() (types.Providers, error) {
		return nil, err
	}
}

func buildFetch(out types.Providers) seeker.FetchSeeds {
	return func() (types.Providers, error) {
		return out, nil
	}
}

func buildStoreWithError(err error) seeker.Store {
	return func(regions types.Providers) error {
		return err
	}
}

func buildStore() seeker.Store {
	return func(pr types.Providers) error {
		return nil
	}
}
