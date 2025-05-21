package cli_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	cli "github.com/kyma-project/gardener-syncer/internal"
	"github.com/stretchr/testify/require"
)

var defaultConfig = cli.Config{
	Gardener: cli.Gardener{
		KubeconfigPath: cli.FlagDefaultGardenerKubeconfigPath,
	},
}

func TestNewConfigFromFlags(t *testing.T) {

	testCases := []struct {
		name          string
		args          []string
		expectedError error
		expectedCfg   cli.Config
	}{
		{
			name:        "OK1: defaults",
			args:        []string{},
			expectedCfg: defaultConfig,
		},
		{
			name: "OK2: kcp-kubeconfig-path override",
			args: []string{
				fmt.Sprintf("-%s", cli.FlagNameGardenerKubeconfigPath), "config.go",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// GIVEN
			os.Args = append([]string{"test"}, testCase.args...)
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			// WHEN
			_, err := cli.NewConfigFromFlags()

			// THEN
			if testCase.expectedError == nil {
				require.NoError(t, err)
			}

			// THEN
			if testCase.expectedError != nil {
				require.ErrorIs(t, err, testCase.expectedError)
				require.ErrorContains(t, err, testCase.expectedError.Error())
			}
		})
	}
}
